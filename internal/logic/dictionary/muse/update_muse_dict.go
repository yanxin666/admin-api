package muse

import (
	"context"
	"muse-admin/internal/model/system"
	"muse-admin/pkg/errs"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMuseDict struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateMuseDict(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMuseDict {
	return &UpdateMuseDict{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateMuseDict) UpdateMuseDict(req *types.UpdateMuseDictReq) (resp *types.Result, err error) {
	err = l.svcCtx.DictWhiteModel.Update(l.ctx, &system.DictWhite{
		Id:     req.Id,
		Key:    req.Key,
		Remark: req.Remark,
	})

	if err != nil {
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}

	return &types.Result{
		Result: true,
	}, nil
}
