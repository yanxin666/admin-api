package muse

import (
	"context"
	"muse-admin/internal/model/system"
	"muse-admin/pkg/errs"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddMuseDict struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddMuseDict(ctx context.Context, svcCtx *svc.ServiceContext) *AddMuseDict {
	return &AddMuseDict{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddMuseDict) AddMuseDict(req *types.AddMuseDictReq) (resp *types.Result, err error) {
	dictWhite, err := l.svcCtx.DictWhiteModel.FindByKey(l.ctx, req.Key)
	if err != nil {
		return nil, err
	}
	if dictWhite != nil {
		return nil, errs.WithCode(err, errs.AddDictionaryErrorCode)
	}

	_, err = l.svcCtx.DictWhiteModel.Insert(l.ctx, &system.DictWhite{
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
