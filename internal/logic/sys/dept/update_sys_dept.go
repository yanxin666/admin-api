package dept

import (
	"context"
	"errors"
	"muse-admin/internal/define"
	"muse-admin/internal/model/workbench"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"
	"slices"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSysDeptLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSysDeptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSysDeptLogic {
	return &UpdateSysDeptLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSysDeptLogic) UpdateSysDept(req *types.UpdateSysDeptReq) error {
	if req.ParentId != define.SysTopParentId {
		_, err := l.svcCtx.SysDeptModel.FindOne(l.ctx, req.ParentId)
		if err != nil {
			return errs.NewCode(errs.ParentDeptIdErrorCode)
		}
	}

	if req.Id == req.ParentId {
		return errs.NewCode(errs.ParentDeptErrorCode)
	}

	dept, err := l.svcCtx.SysDeptModel.FindOneByUniqueKey(l.ctx, req.UniqueKey)
	if !errors.Is(err, workbench.ErrNotFound) && dept.Id != req.Id {
		return errs.NewCode(errs.UpdateDeptUniqueKeyErrorCode)
	}

	deptIds := make([]int64, 0)
	deptIds = l.getSubDept(deptIds, req.Id)
	if slices.Contains(deptIds, req.ParentId) {
		return errs.NewCode(errs.SetParentIdErrorCode)
	}

	sysDept, err := l.svcCtx.SysDeptModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return errs.NewCode(errs.DeptIdErrorCode)
	}

	err = copier.Copy(sysDept, req)
	if err != nil {
		return errs.WithCode(err, errs.ServerErrorCode)
	}

	err = l.svcCtx.SysDeptModel.Update(l.ctx, sysDept)
	if err != nil {
		return errs.WithCode(err, errs.ServerErrorCode)
	}

	return nil
}

func (l *UpdateSysDeptLogic) getSubDept(deptIds []int64, id int64) []int64 {
	deptList, err := l.svcCtx.SysDeptModel.FindSubDept(l.ctx, id)
	if err != nil && err != workbench.ErrNotFound {
		return deptIds
	}

	for _, v := range deptList {
		deptIds = append(deptIds, v.Id)
		deptIds = l.getSubDept(deptIds, v.Id)
	}

	return deptIds
}
