package _import

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/logz"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"muse-admin/internal/define/mqdef"
	taskModel "muse-admin/internal/model/task"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"muse-admin/pkg/excel"
	"strings"
	"time"
)

type Benefit struct {
	svcCtx   *svc.ServiceContext
	taskInfo *taskModel.SyncTask
}

func NewBenefitBuilder(svcCtx *svc.ServiceContext, taskInfo *taskModel.SyncTask) IBuilder {
	return &Benefit{
		svcCtx:   svcCtx,
		taskInfo: taskInfo,
	}
}

func (b *Benefit) GetName() string {
	return "Excel导入用户权益"
}

func (b *Benefit) RepeatSnapshot() bool {
	return false
}

func (b *Benefit) Access() bool {
	return b.taskInfo.Type == 1
}

func (b *Benefit) PreBuild(ctx context.Context, filename string) ([]map[string]string, error) {
	header := map[string]string{
		"所属渠道":         "channel",
		"商品名称":         "name",
		"商品规格":         "spec",
		"权益名称":         "benefit_name",
		"商品数量":         "num",
		"商品金额":         "amount",
		"支付金额(支付的总金额)": "pay_amount",
		"实付金额(支付的现金)":  "real_pay_amount",
		"伴学金支付金额":      "accompany_amount",
		"豆点支付金额":       "bean_amount",
		"支付方式":         "pay_type",
		"订单状态":         "order_status",
		"店铺名称":         "shop_name",
		"直播间(渠道达人)":    "live_name",
		"推广员手机号":       "promoter_phone",
		"渠道订单号":        "channel_order",
		"订单号":          "order_no",
		"手机号":          "phone",
		"支付时间":         "pay_time",
		"退费时间":         "refund_time",
		"退费金额(退费的总金额)": "refund_amount",
		"来源":           "source",
	}
	o := excel.NewExCelImport(
		excel.WithSheetNum(cast.ToInt(b.taskInfo.FileSheet)),
		excel.WithHeader(header),
	)

	return o.GetExcelDataOneSheet(ctx, filename, b.taskInfo.FileName)
}

func (b *Benefit) BuildBefore(_ context.Context, sheetData []map[string]string) (*ExecuteResult, error) {
	resp := &ExecuteResult{
		PreCount:    0,
		SuccessList: make([]BatchContent, 0),
		FailList:    make([]BatchContent, 0),
	}

	for k, v := range sheetData {
		// 先对处理前的原始数据做一个md5的转换，方便后续跟表的md5值匹配查重
		data, _ := json.Marshal(v)
		v["md5"] = util.GenerateMD5Str(string(data))

		// 为空手机号需要过滤
		if v["phone"] == "" {
			resp.FailList = append(resp.FailList, BatchContent{
				Index: k,
				Err:   errors.New("手机号为空，被过滤"),
			})
			continue
		}

		// 删除手机号段内的所有空格
		v["phone"] = strings.ReplaceAll(v["phone"], " ", "")
		// 是否满足手机号格式
		if !util.CheckMobile(v["phone"]) {
			resp.FailList = append(resp.FailList, BatchContent{
				Index: k,
				Err:   errors.New("手机号格式错误，被过滤"),
			})
			continue
		}

		// 过滤部分退款
		if v["order_status"] == "部分退款" {
			resp.FailList = append(resp.FailList, BatchContent{
				Index: k,
				Err:   errors.New("该数据状态为部分退款，被过滤"),
			})
			continue
		}

		// 订单状态转换
		orderStatus, ok := orderStatusMap[v["order_status"]]
		if !ok {
			resp.FailList = append(resp.FailList, BatchContent{
				Index: k,
				Err:   errors.New(fmt.Sprintf("订单状态:%s被过滤，原因：枚举转换失败", v["order_status"])),
			})
			continue
		}
		v["order_status"] = orderStatus

		// 来源转换
		source, ok := sourceMap[v["source"]]
		if !ok {
			resp.FailList = append(resp.FailList, BatchContent{
				Index: k,
				Err:   errors.New(fmt.Sprintf("来源:%s被过滤，原因：枚举转换失败", v["source"])),
			})
			continue
		}
		v["source"] = source

		// 订单号为空时，需要创建一个订单号
		if v["order_no"] == "" {
			orderNo := b.svcCtx.Snowflake.Generate()
			v["order_no"] = cast.ToString(orderNo)
		}
	}

	return resp, nil
}

func (b *Benefit) BuildData(ctx context.Context, sheetData []map[string]string, resp *ExecuteResult) {
	for k, v := range sheetData {
		// 执行量
		_ = b.setRedisCount(ctx, b.taskInfo.Id)

		// 过滤为空的数据源
		if v == nil {
			continue
		}

		// 已成功的数据不重复执行
		detail, _ := b.svcCtx.SyncTaskLogModel.FindOneByMd5AndStatus(ctx, v["md5"], StatusType.Success)
		if detail != nil {
			resp.FailList = append(resp.FailList, BatchContent{
				Index: k,
				Err:   errors.New(fmt.Sprintf("该手机号:%s被过滤，原因：已在「任务%d」中被执行", v["phone"], detail.TaskId)),
			})
			continue
		}

		// 预处理数据记录
		resp.PreCount += 1

		// 用户权益预导入处理
		_, err := b.svcCtx.MQProducer.SendSync(ctx, mqdef.TopicUserBenefitPreImport, types.ExcelImportBenefitModel{
			BenefitName:  v["benefit_name"],
			Phone:        v["phone"],
			ChannelOrder: v["channel_order"],
			OrderNo:      v["order_no"],
			Source:       cast.ToInt64(v["source"]),
			OrderStatus:  cast.ToInt64(v["order_status"]),
		})
		if err != nil {
			// 添加到failList后，继续下一个事务操作
			resp.FailList = append(resp.FailList, BatchContent{
				Index: k,
				Err:   err,
			})
			continue
		}

		// 执行成功
		resp.SuccessList = append(resp.SuccessList, BatchContent{
			Index: k,
		})

		time.Sleep(200 * time.Millisecond) // 设置传输间隔为200毫秒，1秒内发送5条，1分内发送300个
	}
}

func (b *Benefit) BuildAfter(context.Context, *ExecuteResult) error {
	return nil
}

// Redis设置执行量
func (b *Benefit) setRedisCount(ctx context.Context, taskId int64) error {
	key := fmt.Sprintf(RedisKeyImportCount, taskId)
	value, err := b.svcCtx.RedisClient.Get(ctx, key)
	if err != nil {
		return err
	}

	// 若限制不存在就设置
	if value == "" {
		// 过期时间为24小时
		_, err = b.svcCtx.RedisClient.SetnxEx(ctx, key, "1", 86400)
		if err != nil {
			logz.Errorf(ctx, "RedisKeySetNxEx:key:%s,err:%s", key, err)
			return err
		}
	} else {
		_, err = b.svcCtx.RedisClient.Incr(ctx, key)
		if err != nil {
			logz.Errorf(ctx, "RedisKeySetIncr:key:%s,err:%s", key, err)
			return err
		}
	}

	return nil
}
