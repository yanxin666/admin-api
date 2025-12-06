package member

import (
	"context"
	"muse-admin/pkg/errs"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type WhiteDelete struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWhiteDelete(ctx context.Context, svcCtx *svc.ServiceContext) *WhiteDelete {
	return &WhiteDelete{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WhiteDelete) WhiteDelete(req *types.WhiteDeleteReq) (resp *types.Result, err error) {
	if err = l.checkParams(req); err != nil {
		return nil, err
	}

	err = l.svcCtx.UserWhiteModel.Delete(l.ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &types.Result{
		Result: true,
	}, nil
}

// 校验参数
func (l *WhiteDelete) checkParams(req *types.WhiteDeleteReq) error {
	if req.Id <= 0 {
		return errs.NewMsg(errs.ErrCodeParamsAbnormal, "ID不存在").ShowMsg()
	}

	return nil
}
