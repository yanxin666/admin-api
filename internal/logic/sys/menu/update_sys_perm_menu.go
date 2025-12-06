package menu

import (
	"context"
	"encoding/json"
	"errors"
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

type UpdateSysPermMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSysPermMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSysPermMenuLogic {
	return &UpdateSysPermMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSysPermMenuLogic) UpdateSysPermMenu(req *types.UpdateSysPermMenuReq) error {
	userId := tools.GetUserIdByCtx(l.ctx)
	if userId != define.SysSuperUserId {
		for _, v := range req.Perms {
			is, err := l.svcCtx.Redis.Sismember(define.SysPermMenuCachePrefix+strconv.FormatInt(userId, 10), define.SysPermMenuPrefix+v)
			if err != nil || is != true {
				return errs.NewCode(errs.NotPermMenuErrorCode)
			}
		}
	}

	if req.ParentId != define.SysTopParentId {
		parentPermMenu, err := l.svcCtx.SysPermMenuModel.FindOne(l.ctx, req.ParentId)
		if err != nil {
			return errs.NewCode(errs.ParentPermMenuIdErrorCode)
		}

		if parentPermMenu.Type == 2 {
			return errs.NewCode(errs.SetParentTypeErrorCode)
		}
	}

	if req.Id <= define.SysProtectPermMenuMaxId {
		return errs.NewCode(errs.ForbiddenErrorCode)
	}

	if req.Id == req.ParentId {
		return errs.NewCode(errs.ParentPermMenuErrorCode)
	}

	permMenuIds := make([]int64, 0)
	permMenuIds = l.getSubPermMenu(permMenuIds, req.Id)
	if slices.Contains(permMenuIds, req.ParentId) {
		return errs.NewCode(errs.SetParentIdErrorCode)
	}

	permMenu, err := l.svcCtx.SysPermMenuModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return errs.NewCode(errs.PermMenuIdErrorCode)
	}

	err = copier.Copy(permMenu, req)
	if err != nil {
		return errs.WithCode(err, errs.ServerErrorCode)
	}

	bytes, err := json.Marshal(req.Perms)
	if err != nil {
		return errs.WithCode(err, errs.ServerErrorCode)
	}

	permMenu.Perms = string(bytes)
	err = l.svcCtx.SysPermMenuModel.Update(l.ctx, permMenu)
	if err != nil {
		return errs.WithCode(err, errs.ServerErrorCode)
	}

	return nil
}

func (l *UpdateSysPermMenuLogic) getSubPermMenu(permMenuIds []int64, id int64) []int64 {
	permMenuList, err := l.svcCtx.SysPermMenuModel.FindSubPermMenu(l.ctx, id)
	if err != nil && !errors.Is(err, workbench.ErrNotFound) {
		return permMenuIds
	}

	for _, v := range permMenuList {
		permMenuIds = append(permMenuIds, v.Id)
		permMenuIds = l.getSubPermMenu(permMenuIds, v.Id)
	}

	return permMenuIds
}
