package dict

import (
	"context"
	"muse-admin/internal/define"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteConfigDictLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteConfigDictLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteConfigDictLogic {
	return &DeleteConfigDictLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteConfigDictLogic) DeleteConfigDict(req *types.DeleteConfigDictReq) error {
	if req.Id <= define.SysProtectDictionaryMaxId {
		return errs.NewCode(errs.ForbiddenErrorCode)
	}

	total, err := l.svcCtx.SysDictionaryModel.FindCountByParentId(l.ctx, req.Id)
	if err != nil {
		return errs.WithCode(err, errs.ServerErrorCode)
	}

	if total > 0 {
		return errs.NewCode(errs.DeleteDictionaryErrorCode)
	}

	err = l.svcCtx.SysDictionaryModel.Delete(l.ctx, req.Id)
	if err != nil {
		return errs.WithCode(err, errs.ServerErrorCode)
	}

	return nil
}
