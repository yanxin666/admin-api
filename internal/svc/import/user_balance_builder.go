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
	"strconv"
	"strings"
	"time"
)

type UserBalance struct {
	svcCtx   *svc.ServiceContext
	taskInfo *taskModel.SyncTask
}

func NewUserBalanceBuilder(svcCtx *svc.ServiceContext, taskInfo *taskModel.SyncTask) IBuilder {
	return &UserBalance{
		svcCtx:   svcCtx,
		taskInfo: taskInfo,
	}
}

func (s *UserBalance) GetName() string {
	return "Excel导入用户余额"
}

func (s *UserBalance) RepeatSnapshot() bool {
	return false
}

func (s *UserBalance) Access() bool {
	return s.taskInfo.Type == 2
}

func (s *UserBalance) PreBuild(ctx context.Context, filename string) ([]map[string]string, error) {
	header := map[string]string{
		"手机号":       "phone",
		"金额（悟性石数量）": "amount",
		"订单号":       "order_no",
	}
	o := excel.NewExCelImport(
		excel.WithSheetNum(cast.ToInt(s.taskInfo.FileSheet)),
		excel.WithHeader(header),
	)

	return o.GetExcelDataOneSheet(ctx, filename, s.taskInfo.FileName)
}

func (s *UserBalance) BuildBefore(_ context.Context, sheetData []map[string]string) (*ExecuteResult, error) {
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

		// 订单号为空需要过滤
		if v["order_no"] == "" {
			resp.FailList = append(resp.FailList, BatchContent{
				Index: k,
				Err:   errors.New("订单号为空，被过滤"),
			})
			continue
		}

		// 数量为空需要过滤
		if v["amount"] == "" {
			resp.FailList = append(resp.FailList, BatchContent{
				Index: k,
				Err:   errors.New("数量为空，被过滤"),
			})
			continue
		}
		if !isIntegerAndGtZero(v["amount"]) {
			resp.FailList = append(resp.FailList, BatchContent{
				Index: k,
				Err:   errors.New("数量为非大于0的整型，被过滤"),
			})
			continue
		}
	}

	return resp, nil
}

func (s *UserBalance) BuildData(ctx context.Context, sheetData []map[string]string, resp *ExecuteResult) {
	for k, v := range sheetData {
		// 过滤为空的数据源
		if v == nil {
			continue
		}

		// 已成功的数据不重复执行
		detail, _ := s.svcCtx.SyncTaskLogModel.FindOneByMd5AndStatus(ctx, v["md5"], StatusType.Success)
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
		_, err := s.svcCtx.MQProducer.SendSync(ctx, mqdef.TopicUserBalanceImport, types.ExcelImportBalanceData{
			Phone:   v["phone"],
			Amount:  cast.ToInt64(v["amount"]),
			OrderNo: v["order_no"],
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

		time.Sleep(30 * time.Millisecond) // 设置传输间隔为30毫秒，1秒内发送30条
	}
}

func (s *UserBalance) BuildAfter(context.Context, *ExecuteResult) error {
	return nil
}

// Redis设置执行量
func (s *UserBalance) setRedisCount(ctx context.Context, taskId int64) error {
	key := fmt.Sprintf(RedisKeyImportCount, taskId)
	value, err := s.svcCtx.RedisClient.Get(ctx, key)
	if err != nil {
		return err
	}

	// 若限制不存在就设置
	if value == "" {
		// 过期时间为每天23:59:59
		_, err = s.svcCtx.RedisClient.SetnxEx(ctx, key, "1", 86400)
		if err != nil {
			logz.Errorf(ctx, "RedisKeySetNxEx:key:%s,err:%s", key, err)
			return err
		}
	} else {
		_, err = s.svcCtx.RedisClient.Incr(ctx, key)
		if err != nil {
			logz.Errorf(ctx, "RedisKeySetIncr:key:%s,err:%s", key, err)
			return err
		}
	}

	return nil
}

// 判断一个字符串是否可以转换为整型且大于0
func isIntegerAndGtZero(s string) bool {
	value, err := strconv.Atoi(s)
	if err != nil {
		return false // 转换失败，说明字符串不是有效的整型
	}
	return value > 0 // 判断转换后的值是否大于0
}

// 判断一个字符串是否可以转换为整型
func isInteger(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
