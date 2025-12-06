package dept

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

type AddSysDeptLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddSysDeptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddSysDeptLogic {
	return &AddSysDeptLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddSysDeptLogic) AddSysDept(req *types.AddSysDeptReq) error {
	_, err := l.svcCtx.SysDeptModel.FindOneByUniqueKey(l.ctx, req.UniqueKey)
	if errors.Is(err, workbench.ErrNotFound) {
		if req.ParentId != define.SysTopParentId {
			_, err := l.svcCtx.SysDeptModel.FindOne(l.ctx, req.ParentId)
			if err != nil {
				return errs.NewCode(errs.ParentDeptIdErrorCode)
			}
		}

		var sysDept = new(workbench.Dept)
		err = copier.Copy(sysDept, req)
		if err != nil {
			return errs.WithCode(err, errs.ServerErrorCode)
		}
		_, err = l.svcCtx.SysDeptModel.Insert(l.ctx, sysDept)
		if err != nil {
			return errs.WithCode(err, errs.ServerErrorCode)
		}
		return nil
	} else {
		return errs.NewCode(errs.AddDeptErrorCode)
	}
}
