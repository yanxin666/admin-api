package supervisor

import (
	"context"
	"database/sql"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"muse-admin/internal/model/supervisor"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type InteractionUpdate struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInteractionUpdate(ctx context.Context, svcCtx *svc.ServiceContext) *InteractionUpdate {
	return &InteractionUpdate{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InteractionUpdate) InteractionUpdate(req *types.InteractionUpdateReq) error {
	if err := l.checkParams(req); err != nil {
		return err
	}
	_, err := l.svcCtx.InteractionModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return errs.NewMsg(errs.ErrCodeProgram, "互动数据不存在")
	}

	_, err = l.svcCtx.InteractionModel.UpdateFillFieldsById(l.ctx, req.Id, &supervisor.Interaction{
		Name:        req.Name,
		Description: req.Description,
		TeachingData: sql.NullString{
			String: req.TeachingData,
			Valid:  req.TeachingData != "",
		},
		Data: sql.NullString{
			String: req.Data,
			Valid:  req.Data != "",
		},
	})
	if err != nil {
		return err
	}

	return nil
}

// 校验参数
func (l *InteractionUpdate) checkParams(req *types.InteractionUpdateReq) error {
	if req.TeachingData != "" && !util.IsJSON(req.TeachingData) {
		return errs.NewMsg(errs.ErrCodeProgram, "教学数据格式不正确")
	}

	if req.Data != "" && !util.IsJSON(req.Data) {
		return errs.NewMsg(errs.ErrCodeProgram, "互动数据格式不正确")
	}
	return nil
}
