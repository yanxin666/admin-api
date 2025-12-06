package user

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"encoding/json"
	"muse-admin/internal/define"
	"muse-admin/internal/svc"
	"muse-admin/internal/tools"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"
	"slices"
	"strconv"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSysUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSysUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSysUserLogic {
	return &UpdateSysUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSysUserLogic) UpdateSysUser(req *types.UpdateSysUserReq) error {
	currentUserId := tools.GetUserIdByCtx(l.ctx)
	var currentUserRoleIds []int64
	var roleIds []int64
	if currentUserId == define.SysSuperUserId {
		sysRoleList, _ := l.svcCtx.SysRoleModel.FindAll(l.ctx)
		for _, role := range sysRoleList {
			currentUserRoleIds = append(currentUserRoleIds, role.Id)
			roleIds = append(roleIds, role.Id)
		}

	} else {
		currentUser, _ := l.svcCtx.SysUserModel.FindOne(l.ctx, currentUserId)
		err := json.Unmarshal([]byte(currentUser.RoleIds), &currentUserRoleIds)
		if err != nil {
			return errs.WithCode(err, errs.ServerErrorCode)
		}

		roleIds = append(roleIds, currentUserRoleIds...)
	}

	editUser, err := l.svcCtx.SysUserModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return errs.NewCode(errs.UserIdErrorCode)
	}

	var editUserRoleIds []int64
	err = json.Unmarshal([]byte(editUser.RoleIds), &editUserRoleIds)
	if err != nil {
		return errs.WithCode(err, errs.ServerErrorCode)
	}
	roleIds = append(roleIds, editUserRoleIds...)

	for _, id := range req.RoleIds {
		if !slices.Contains(roleIds, id) {
			return errs.NewCode(errs.AssigningRolesErrorCode)
		}
	}

	for _, id := range util.ArrayDifference(editUserRoleIds, currentUserRoleIds) {
		if !slices.Contains(req.RoleIds, id) {
			return errs.NewCode(errs.AssigningRolesErrorCode)
		}
	}

	_, err = l.svcCtx.SysDeptModel.FindOne(l.ctx, req.DeptId)
	if err != nil {
		return errs.NewCode(errs.DeptIdErrorCode)
	}

	_, err = l.svcCtx.SysProfessionModel.FindOne(l.ctx, req.ProfessionId)
	if err != nil {
		return errs.NewCode(errs.ProfessionIdErrorCode)
	}

	_, err = l.svcCtx.SysJobModel.FindOne(l.ctx, req.JobId)
	if err != nil {
		return errs.NewCode(errs.JobIdErrorCode)
	}

	err = copier.Copy(editUser, req)
	if err != nil {
		return errs.WithCode(err, errs.ServerErrorCode)
	}

	bytes, err := json.Marshal(req.RoleIds)
	if err != nil {
		return errs.WithCode(err, errs.ServerErrorCode)
	}

	_, err = l.svcCtx.Redis.Del(define.SysPermMenuCachePrefix + strconv.FormatInt(editUser.Id, 10))
	_, err = l.svcCtx.Redis.Del(define.SysOnlineUserCachePrefix + strconv.FormatInt(editUser.Id, 10))
	editUser.RoleIds = string(bytes)
	err = l.svcCtx.SysUserModel.Update(l.ctx, editUser)
	if err != nil {
		return errs.WithCode(err, errs.ServerErrorCode)
	}

	return nil
}
