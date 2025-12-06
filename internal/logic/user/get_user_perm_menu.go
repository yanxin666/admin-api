package user

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"encoding/json"
	"muse-admin/internal/define"
	"muse-admin/internal/model/workbench"
	"muse-admin/internal/svc"
	"muse-admin/internal/tools"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"
	"slices"
	"strconv"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserPermMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserPermMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserPermMenuLogic {
	return &GetUserPermMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserPermMenuLogic) GetUserPermMenu() (resp *types.UserPermMenuResp, err error) {
	userId := tools.GetUserIdByCtx(l.ctx)

	online, err := l.svcCtx.Redis.Get(define.SysOnlineUserCachePrefix + strconv.FormatInt(userId, 10))
	if err != nil || online == "" {
		return nil, errs.NewCode(errs.AuthErrorCode)
	}

	// 查询用户信息
	user, err := l.svcCtx.SysUserModel.FindOne(l.ctx, userId)
	if err != nil {
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}

	var roles []int64
	// 用户所属角色
	err = json.Unmarshal([]byte(user.RoleIds), &roles)
	if err != nil {
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}

	var permMenu []int64
	var userPermMenu []*workbench.PermMenu

	userPermMenu, permMenu, err = l.countUserPermMenu(roles, permMenu)
	if err != nil {
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}

	var menu types.Menu
	menuList := make([]types.Menu, 0)
	permList := make([]string, 0)
	_, err = l.svcCtx.Redis.Del(define.SysPermMenuCachePrefix + strconv.FormatInt(userId, 10))
	for _, v := range userPermMenu {
		err := copier.Copy(&menu, &v)
		if err != nil {
			return nil, errs.WithCode(err, errs.ServerErrorCode)
		}

		if menu.Type != define.SysDefaultPermType {
			menuList = append(menuList, menu)
		}
		var permArray []string
		err = json.Unmarshal([]byte(v.Perms), &permArray)
		if err != nil {
			return nil, errs.WithCode(err, errs.ServerErrorCode)
		}

		for _, p := range permArray {
			_, err := l.svcCtx.Redis.Sadd(define.SysPermMenuCachePrefix+strconv.FormatInt(userId, 10), define.SysPermMenuPrefix+p)
			if err != nil {
				return nil, errs.WithCode(err, errs.ServerErrorCode)
			}
			permList = append(permList, "/"+p)
		}

	}

	return &types.UserPermMenuResp{Menus: menuList, Perms: util.ArrayUniqueValue[string](permList)}, nil
}

func (l *GetUserPermMenuLogic) countUserPermMenu(roles []int64, permMenu []int64) ([]*workbench.PermMenu, []int64, error) {
	if slices.Contains(roles, define.SysSuperRoleId) {
		sysPermMenus, err := l.svcCtx.SysPermMenuModel.FindAll(l.ctx)
		if err != nil {
			return nil, permMenu, err
		}

		return sysPermMenus, permMenu, nil
	} else {
		for _, roleId := range roles {
			// 查询角色信息
			role, err := l.svcCtx.SysRoleModel.FindOne(l.ctx, roleId)
			if err != nil && err != workbench.ErrNotFound {
				return nil, permMenu, errs.WithCode(err, errs.ServerErrorCode)
			}

			var perms []int64
			// 角色所拥有的权限id
			err = json.Unmarshal([]byte(role.PermMenuIds), &perms)
			if err != nil {
				return nil, permMenu, errs.WithCode(err, errs.ServerErrorCode)
			}

			if role.Status != 0 {
				permMenu = append(permMenu, perms...)
			}
			// 汇总用户所属角色权限id
			permMenu = l.getRolePermMenu(permMenu, roleId)
		}

		// 过滤重复的权限id
		permMenu = util.ArrayUniqueValue[int64](permMenu)
		var ids string
		for i, id := range permMenu {
			if i == 0 {
				ids = strconv.FormatInt(id, 10)
				continue
			}
			ids = ids + "," + strconv.FormatInt(id, 10)
		}

		if len(ids) == 0 {
			return nil, permMenu, nil
		}

		// 根据权限id获取具体权限
		sysPermMenus, err := l.svcCtx.SysPermMenuModel.FindByIds(l.ctx, ids)
		if err != nil {
			return nil, permMenu, err
		}

		return sysPermMenus, permMenu, nil
	}
}

func (l *GetUserPermMenuLogic) getRolePermMenu(perms []int64, roleId int64) []int64 {
	roles, err := l.svcCtx.SysRoleModel.FindSubRole(l.ctx, roleId)
	if err != nil && err != workbench.ErrNotFound {
		return perms
	}

	for _, role := range roles {
		var subPerms []int64
		err = json.Unmarshal([]byte(role.PermMenuIds), &subPerms)
		perms = append(perms, subPerms...)
		perms = l.getRolePermMenu(perms, role.Id)
	}

	return perms
}
