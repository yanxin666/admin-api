package menu

import (
	"context"
	"encoding/json"
	"muse-admin/internal/define"
	"muse-admin/internal/svc"
	"muse-admin/internal/tools"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"
	"slices"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSysPermMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteSysPermMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSysPermMenuLogic {
	return &DeleteSysPermMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSysPermMenuLogic) DeleteSysPermMenu(req *types.DeleteSysPermMenuReq) error {
	currentUserId := tools.GetUserIdByCtx(l.ctx)
	if currentUserId != define.SysSuperUserId {
		var currentUserPermMenuIds []int64
		currentUserPermMenuIds = l.getCurrentUserPermMenuIds(currentUserId)
		if !slices.Contains(currentUserPermMenuIds, req.Id) {
			return errs.NewCode(errs.NotPermMenuErrorCode)
		}
	}

	if req.Id <= define.SysProtectPermMenuMaxId {
		return errs.NewCode(errs.ForbiddenErrorCode)
	}

	count, _ := l.svcCtx.SysPermMenuModel.FindCountByParentId(l.ctx, req.Id)
	if count != 0 {
		return errs.NewCode(errs.DeletePermMenuErrorCode)
	}

	err := l.svcCtx.SysPermMenuModel.Delete(l.ctx, req.Id)
	if err != nil {
		return errs.WithCode(err, errs.ServerErrorCode)
	}

	return nil
}

func (l *DeleteSysPermMenuLogic) getCurrentUserPermMenuIds(currentUserId int64) (ids []int64) {
	var currentPermMenuIds []int64
	if currentUserId != define.SysSuperUserId {
		var currentUserRoleIds []int64
		var roleIds []int64
		currentUser, _ := l.svcCtx.SysUserModel.FindOne(l.ctx, currentUserId)
		_ = json.Unmarshal([]byte(currentUser.RoleIds), &currentUserRoleIds)
		roleIds = append(roleIds, currentUserRoleIds...)
		var ids string
		for i, v := range roleIds {
			if i == 0 {
				ids = strconv.FormatInt(v, 10)
			}
			ids = ids + "," + strconv.FormatInt(v, 10)
		}

		sysRoles, _ := l.svcCtx.SysRoleModel.FindByIds(l.ctx, ids)
		var rolePermMenus []int64
		for _, v := range sysRoles {
			err := json.Unmarshal([]byte(v.PermMenuIds), &rolePermMenus)
			if err != nil {
				return nil
			}
			currentPermMenuIds = append(currentPermMenuIds, rolePermMenus...)
		}
	}

	return currentPermMenuIds
}
