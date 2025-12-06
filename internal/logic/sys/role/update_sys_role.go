package role

import (
	"context"
	"encoding/json"
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

type UpdateSysRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSysRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSysRoleLogic {
	return &UpdateSysRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSysRoleLogic) UpdateSysRole(req *types.UpdateSysRoleReq) error {
	if req.ParentId != define.SysTopParentId {
		_, err := l.svcCtx.SysRoleModel.FindOne(l.ctx, req.ParentId)
		if err != nil {
			return errs.NewCode(errs.ParentRoleIdErrorCode)
		}
	}

	if req.Id == define.SysSuperRoleId {
		return errs.NewCode(errs.NotPermMenuErrorCode)
	}

	if req.Id == req.ParentId {
		return errs.NewCode(errs.ParentRoleErrorCode)
	}

	role, err := l.svcCtx.SysRoleModel.FindOneByUniqueKey(l.ctx, req.UniqueKey)
	if !errors.Is(err, workbench.ErrNotFound) && role.Id != req.Id {
		return errs.NewCode(errs.UpdateRoleUniqueKeyErrorCode)
	}

	roleIds := make([]int64, 0)
	roleIds = l.getSubRole(roleIds, req.Id)
	if slices.Contains(roleIds, req.ParentId) {
		return errs.NewCode(errs.SetParentIdErrorCode)
	}

	sysRole, err := l.svcCtx.SysRoleModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return errs.NewCode(errs.RoleIdErrorCode)
	}

	err = copier.Copy(sysRole, req)
	if err != nil {
		return errs.WithCode(err, errs.ServerErrorCode)
	}

	bytes, err := json.Marshal(req.PermMenuIds)
	if err != nil {
		return errs.WithCode(err, errs.ServerErrorCode)
	}

	sysRole.PermMenuIds = string(bytes)
	err = l.svcCtx.SysRoleModel.Update(l.ctx, sysRole)
	if err != nil {
		return errs.WithCode(err, errs.ServerErrorCode)
	}

	return nil
}

func (l *UpdateSysRoleLogic) getSubRole(roleIds []int64, id int64) []int64 {
	roleList, err := l.svcCtx.SysRoleModel.FindSubRole(l.ctx, id)
	if err != nil && !errors.Is(err, workbench.ErrNotFound) {
		return roleIds
	}

	for _, v := range roleList {
		roleIds = append(roleIds, v.Id)
		roleIds = l.getSubRole(roleIds, v.Id)
	}

	return roleIds
}
