package dict

import (
	"context"
	"errors"
	"muse-admin/internal/define"
	"muse-admin/internal/model/workbench"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddConfigDictLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddConfigDictLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddConfigDictLogic {
	return &AddConfigDictLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddConfigDictLogic) AddConfigDict(req *types.AddConfigDictReq) error {
	if req.ParentId != define.SysTopParentId {
		_, err := l.svcCtx.SysDictionaryModel.FindOne(l.ctx, req.ParentId)
		if err != nil {
			return errs.WithCode(err, errs.ParentDictionaryIdErrorCode)
		}
	}
	_, err := l.svcCtx.SysDictionaryModel.FindOneByUniqueKey(l.ctx, req.UniqueKey)
	if errors.Is(err, workbench.ErrNotFound) {
		var dictionary = new(workbench.Dictionary)
		err = copier.Copy(dictionary, req)
		if err != nil {
			return errs.WithCode(err, errs.ServerErrorCode)
		}
		_, err = l.svcCtx.SysDictionaryModel.Insert(l.ctx, dictionary)
		if err != nil {
			return errs.WithCode(err, errs.ServerErrorCode)
		}

		return nil
	} else {

		return errs.NewCode(errs.AddDictionaryErrorCode)
	}
}
