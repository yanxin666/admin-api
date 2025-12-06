package dict

import (
	"context"
	"muse-admin/internal/define"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateConfigDictLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateConfigDictLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateConfigDictLogic {
	return &UpdateConfigDictLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateConfigDictLogic) UpdateConfigDict(req *types.UpdateConfigDictReq) error {
	if req.ParentId != define.SysTopParentId {
		_, err := l.svcCtx.SysDictionaryModel.FindOne(l.ctx, req.ParentId)
		if err != nil {
			return errs.WithCode(err, errs.ParentDictionaryIdErrorCode)
		}
	}

	configDictionary, err := l.svcCtx.SysDictionaryModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return errs.WithCode(err, errs.DictionaryIdErrorCode)
	}

	err = copier.Copy(configDictionary, req)
	if err != nil {
		return errs.WithCode(err, errs.ServerErrorCode)
	}

	err = l.svcCtx.SysDictionaryModel.Update(l.ctx, configDictionary)
	if err != nil {
		return errs.WithCode(err, errs.ServerErrorCode)
	}

	return nil
}
