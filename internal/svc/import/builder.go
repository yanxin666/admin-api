package _import

import (
	"context"
	"database/sql"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logc"
	taskModel "muse-admin/internal/model/task"
	"muse-admin/internal/svc"
	"muse-admin/internal/svc/public"
	"muse-admin/pkg/errs"
)

type ExecuteResult struct {
	PreCount    int64          `json:"preCount"`    // 预处理条数
	SuccessList []BatchContent `json:"successList"` // 添加成功的题目集合
	FailList    []BatchContent `json:"failList"`    // 添加失败的题目集合
}
type BatchContent struct {
	Index int   // 第几条数据
	Err   error // 错误信息
}

type IBuilder interface {
	GetName() string                                                                        // 获取名称
	Access() bool                                                                           // builder准入
	PreBuild(ctx context.Context, filename string) ([]map[string]string, error)             // 预构建数据
	RepeatSnapshot() bool                                                                   // 是否记录重复的快照
	BuildBefore(ctx context.Context, sheetData []map[string]string) (*ExecuteResult, error) // 处理数据的前置准备操作 如：数据转换等动作
	BuildData(context.Context, []map[string]string, *ExecuteResult)                         // 开始处理数据
	BuildAfter(context.Context, *ExecuteResult) error                                       // 处理数据的后置收尾操作
	setRedisCount(ctx context.Context, taskId int64) error                                  // 处理时的进度
	// ... 添加更多动作
}

type BaseBuilder struct {
	svcCtx           *svc.ServiceContext
	taskInfo         *taskModel.SyncTask
	downloadFilename string
	sheetData        []map[string]string // 所选中的表对象数据
	alarmService     *public.BizAlarm
}

func NewBuilder(svcCtx *svc.ServiceContext, taskInfo *taskModel.SyncTask, filename string) *BaseBuilder {
	return &BaseBuilder{
		svcCtx:           svcCtx,
		taskInfo:         taskInfo,
		downloadFilename: filename,
		alarmService:     public.NewAlarmService(svcCtx),
	}
}

// BuildStage 构建导入各种场景所需要的阶段
func (b *BaseBuilder) BuildStage(ctx context.Context) {
	var (
		resp        *ExecuteResult
		err         error
		builderList []IBuilder
		StageList   []func(context.Context) IBuilder
	)

	// stage注册器，所有相关的stage全部在这里归纳
	StageList = []func(context.Context) IBuilder{
		// b.getGroupBuilder,   // 组题数据导入
		b.getBenefitBuilder,     // 权益数据导入
		b.getUserBalanceBuilder, // 用户悟性石数据导入
		// ... todo 其余导入逻辑接入
	}

	for _, layoutItem := range StageList {
		if iBuilder := layoutItem(ctx); iBuilder != nil {
			builderList = append(builderList, iBuilder)
		}
	}

	for _, iBuilder := range builderList {
		if !iBuilder.Access() {
			continue
		}
		logc.Infof(ctx, "%s已准入,type:%d", iBuilder.GetName(), b.taskInfo.Type)

		// 预构建数据
		b.sheetData, err = iBuilder.PreBuild(ctx, b.downloadFilename)
		if err != nil {
			logc.Errorf(ctx, "%s构建数据异常,err:%v", iBuilder.GetName(), err)
			_ = b.alarmService.ImportAlarm(ctx, b.taskInfo.OperateId, "导入数据-预构建数据错误", err)
			_ = b.svcCtx.SyncTaskModel.UpdateFailById(ctx, b.taskInfo.Id, "文件内容读取失败，请检查更正后重试")
			continue
		}

		// 数据快照
		if err = b.snapshot(ctx, iBuilder.RepeatSnapshot()); err != nil {
			_ = b.alarmService.ImportAlarm(ctx, b.taskInfo.OperateId, "导入数据-数据快照记录错误", err)
			continue
		}

		// 记录开始时间以及处理总数
		if err = b.begin(ctx); err != nil {
			_ = b.alarmService.ImportAlarm(ctx, b.taskInfo.OperateId, "导入数据-预处理更改错误", err)
			continue
		}

		// 执行任务前置操作
		resp, err = iBuilder.BuildBefore(ctx, b.sheetData)
		if err != nil {
			_ = b.alarmService.ImportAlarm(ctx, b.taskInfo.OperateId, "导入数据-前置操作错误", err)
			continue
		}

		// 对前置操作时不符合条件的元素进行空赋值，确保后续的数据出错时有正确的下标
		if len(resp.FailList) != 0 {
			for _, v := range resp.FailList {
				b.sheetData[v.Index] = nil
			}
		}

		// 执行任务处理
		iBuilder.BuildData(ctx, b.sheetData, resp)

		// 执行任务后置操作
		if err = iBuilder.BuildAfter(ctx, resp); err != nil {
			_ = b.alarmService.ImportAlarm(ctx, b.taskInfo.OperateId, "导入数据-后置收尾错误", err)
			continue
		}

		// 记录结果
		if err = b.result(ctx, resp); err != nil {
			_ = b.alarmService.ImportAlarm(ctx, b.taskInfo.OperateId, "导入数据-记录结果失败,请相关人员检查", err)
			continue
		}
	}
}

// 任务开始时需要更新的前置操作
func (b *BaseBuilder) begin(ctx context.Context) error {
	// 记录更新时间以及处理总数
	_, err := b.svcCtx.SyncTaskModel.UpdateFillFieldsById(ctx, b.taskInfo.Id, &taskModel.SyncTask{
		StartTime: util.GetStandardNowDatetime(),
		Total:     cast.ToInt64(len(b.sheetData)),
		Status:    StatusType.Going,
	})
	if err != nil {
		return errs.WithMsg(err, errs.ErrCodeAbnormal, "更新前置动作失败")
	}

	return nil
}

// 前置操作结束后的数据快照记录
func (b *BaseBuilder) snapshot(ctx context.Context, repeatSnapshot bool) error {
	for k, v := range b.sheetData {
		// 所有数据渲染没问题后，做数据快照
		data, _ := json.Marshal(v)
		md5 := util.GenerateMD5Str(string(data))
		detail, err := b.svcCtx.SyncTaskLogModel.FindOneByMd5(ctx, md5)
		if err != nil {
			return errs.WithMsg(err, errs.ErrCodeAbnormal, "查询详情快照失败")
		}

		// 已有数据 && 不记录重复的快照
		if detail != nil && !repeatSnapshot {
			continue
		}

		// 记录快照
		_, err = b.svcCtx.SyncTaskLogModel.Insert(ctx, &taskModel.SyncTaskLog{
			TaskId: b.taskInfo.Id,
			Index:  cast.ToInt64(k),
			Data:   string(data),
			Md5:    md5,
			Status: StatusType.Going,
		})
		if err != nil {
			return errs.WithMsg(err, errs.ErrCodeAbnormal, "数据详情快照失败")
		}
	}

	return nil
}

// 此次任务执行结束后的总结操作
func (b *BaseBuilder) result(ctx context.Context, result *ExecuteResult) error {
	var (
		errMsg    string
		errorMsg  sql.NullString
		status    int64
		isErrLong bool
	)

	errorMsg = sql.NullString{
		String: "暂无",
		Valid:  true,
	}
	status = StatusType.Success

	failLen := len(result.FailList)                                 // 失败总条数「包含过滤条数」
	filterCount := cast.ToInt64(len(b.sheetData)) - result.PreCount // 过滤条数 = 总数 - 预处理条数
	failNum := cast.ToInt64(failLen) - filterCount                  // 真正失败的条数 = 失败总条数 - 过滤条数
	if failLen > 0 {
		// 失败错误太多，去详情中查询即可
		if failLen > 30 {
			isErrLong = true
			errMsg = fmt.Sprintf("部分数据异常条数太多\n过滤数据为：%d条\n失败数据为：%d条", filterCount, failNum)
			if failNum > 0 {
				errMsg += "，请点击详情查询"
			}
		}

		// 更新失败的任务详情内容
		for _, v := range result.FailList {
			if !isErrLong {
				errMsg += fmt.Sprintf("第%d条：%v \n", v.Index+1, v.Err)
			}
			// 过滤不存在的数据
			ok, _ := b.svcCtx.SyncTaskLogModel.FindOneByTaskIdAndIndex(ctx, b.taskInfo.Id, cast.ToInt64(v.Index))
			if ok == nil {
				continue
			}

			// 更新失败原因
			_, _ = b.svcCtx.SyncTaskLogModel.UpdateFillFields(ctx, b.taskInfo.Id, cast.ToInt64(v.Index), &taskModel.SyncTaskLog{
				Status:    StatusType.Fail,
				ErrorsMsg: v.Err.Error(),
			})
		}
		errorMsg = sql.NullString{
			String: errMsg,
			Valid:  true,
		}
		status = StatusType.Fail
	}

	// 更新任务结果
	_, err := b.svcCtx.SyncTaskModel.UpdateFillFieldsById(ctx, b.taskInfo.Id, &taskModel.SyncTask{
		Status:     status,
		EndTime:    util.GetStandardNowDatetime(),
		ErrorMsg:   errorMsg,
		FilterNum:  filterCount,
		PreNum:     result.PreCount,
		SuccessNum: cast.ToInt64(len(result.SuccessList)),
		FailNum:    failNum,
	})
	if err != nil {
		return errs.WithMsg(err, errs.ErrCodeAbnormal, fmt.Sprintf("更新任务结果失败 任务ID：%d", b.taskInfo.Id))
	}

	// 更新成功的任务详情
	if len(result.SuccessList) > 0 {
		_ = b.svcCtx.SyncTaskLogModel.UpdateStatusSuc(ctx, b.taskInfo.Id)
	}

	return nil
}

// func (b *BaseBuilder) getGroupBuilder(context.Context) IBuilder {
// 	return NewGroupBuilder(b.svcCtx, b.taskInfo, b.filename)
// }

func (b *BaseBuilder) getBenefitBuilder(context.Context) IBuilder {
	return NewBenefitBuilder(b.svcCtx, b.taskInfo)
}

func (b *BaseBuilder) getUserBalanceBuilder(context.Context) IBuilder {
	return NewUserBalanceBuilder(b.svcCtx, b.taskInfo)
}
