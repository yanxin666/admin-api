package task

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"muse-admin/internal/model/tools"
	"muse-admin/internal/svc"
	_import "muse-admin/internal/svc/import"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskDetailList struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTaskDetailList(ctx context.Context, svcCtx *svc.ServiceContext) *TaskDetailList {
	return &TaskDetailList{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TaskDetailList) TaskDetailList(req *types.TaskDetailListReq) (resp *types.TaskDetailListResp, err error) {
	resp = &types.TaskDetailListResp{
		List:       make([]types.TaskDetailLog, 0),
		Pagination: types.Pagination{},
	}

	if req.TaskId <= 0 {
		return nil, errs.NewCode(errs.ServerErrorCode)
	}

	taskInfo, err := l.svcCtx.SyncTaskModel.FindOne(l.ctx, req.TaskId)
	if err != nil && !errors.Is(err, sqlc.ErrNotFound) {
		return nil, errs.NewCode(errs.ServerErrorCode)
	}

	if taskInfo == nil {
		return nil, errs.NewMsg(errs.ServerErrorCode, "任务不存在").ShowMsg()
	}

	// 构建查询条件
	condition := tools.FilterConditions(req)
	delete(condition, "keyword")

	// 获取任务详情Log列表
	list, total, err := l.svcCtx.SyncTaskLogModel.FindPageByCondition(l.ctx, req.Page, req.Limit, req.Keyword, condition)
	if err != nil {
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}
	for _, v := range list {
		log := types.TaskDetailLog{
			Index:     v.Index + 1,
			Data:      v.Data,
			Status:    v.Status,
			ErrorsMsg: v.ErrorsMsg,
		}
		resp.List = append(resp.List, log)
	}

	// 分页信息
	resp.Pagination = types.Pagination{
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	}

	key := fmt.Sprintf(_import.RedisKeyImportCount, req.TaskId)
	value, _ := l.svcCtx.RedisClient.Get(l.ctx, key)
	if value != "" {
		resp.Progress = cast.ToInt64(value) / taskInfo.Total * 100
	} else {
		resp.Progress = 100
	}

	return resp, nil
}
