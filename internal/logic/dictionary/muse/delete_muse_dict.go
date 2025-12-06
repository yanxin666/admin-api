package muse

import (
	"context"
	"muse-admin/pkg/errs"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMuseDict struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteMuseDict(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMuseDict {
	return &DeleteMuseDict{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteMuseDict) DeleteMuseDict(req *types.DeleteMuseDictReq) (resp *types.Result, err error) {
	err = l.svcCtx.DictWhiteModel.Delete(l.ctx, req.Id)
	if err != nil {
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}

	return &types.Result{
		Result: true,
	}, nil
}
