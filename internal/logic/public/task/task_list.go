package task

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"fmt"
	"github.com/spf13/cast"
	"math"
	"muse-admin/internal/model/tools"
	_import "muse-admin/internal/svc/import"
	"muse-admin/pkg/errs"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskList struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTaskList(ctx context.Context, svcCtx *svc.ServiceContext) *TaskList {
	return &TaskList{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TaskList) TaskList(req *types.TaskListReq) (resp *types.TaskListResp, err error) {
	resp = &types.TaskListResp{
		List:       make([]types.TaskListInfo, 0),
		Pagination: types.Pagination{},
	}

	// 构建查询条件
	condition := tools.FilterConditions(req)
	// 获取任务列表
	list, total, err := l.svcCtx.SyncTaskModel.FindPageByCondition(l.ctx, req.Page, req.Limit, condition)
	if err != nil {
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}
	if len(list) == 0 {
		return resp, nil
	}

	operateIds, _ := util.ArrayColumn(list, "OperateId")
	operateIds = util.ArrayUniqueValue(operateIds)
	userInfo, err := l.svcCtx.SysUserModel.BatchByUserIds(l.ctx, operateIds)
	if err != nil {
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}

	for _, v := range list {
		var progress int64
		key := fmt.Sprintf(_import.RedisKeyImportCount, v.Id)
		value, _ := l.svcCtx.RedisClient.Get(l.ctx, key)
		if value != "" {
			progress = int64(math.Round(float64(cast.ToInt64(value)) / float64(v.Total) * 100))
		} else {
			progress = 100
		}

		log := types.TaskListInfo{
			Id: v.Id,
			ImportFileInfo: types.ImportFileInfo{
				Filename:      v.FileName,
				FileId:        v.FileId,
				Type:          v.Type,
				FileSheet:     v.FileSheet,
				FileSheetName: v.FileSheetName,
			},
			Status:      v.Status,
			StartTime:   v.StartTime,
			EndTime:     v.EndTime,
			ErrorsMsg:   v.ErrorMsg.String,
			Total:       v.Total,
			FilterNum:   v.FilterNum,
			PreNum:      v.PreNum,
			SuccessNum:  v.SuccessNum,
			FailNum:     v.FailNum,
			Progress:    progress,
			OperateName: userInfo[v.OperateId].Username,
		}
		resp.List = append(resp.List, log)
	}

	// 分页信息
	resp.Pagination = types.Pagination{
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	}
	return resp, nil
}
