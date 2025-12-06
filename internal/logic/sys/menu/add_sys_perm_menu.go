package menu

import (
	"context"
	"encoding/json"
	"muse-admin/internal/define"
	"muse-admin/internal/model/workbench"
	"muse-admin/internal/svc"
	"muse-admin/internal/tools"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"
	"strconv"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddSysPermMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddSysPermMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddSysPermMenuLogic {
	return &AddSysPermMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddSysPermMenuLogic) AddSysPermMenu(req *types.AddSysPermMenuReq) error {
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

	var permMenu = new(workbench.PermMenu)
	err := copier.Copy(permMenu, req)
	if err != nil {
		return errs.WithCode(err, errs.ServerErrorCode)
	}

	bytes, err := json.Marshal(req.Perms)
	if err != nil {
		return errs.WithCode(err, errs.ServerErrorCode)
	}

	permMenu.Perms = string(bytes)
	_, err = l.svcCtx.SysPermMenuModel.Insert(l.ctx, permMenu)
	if err != nil {
		return errs.WithCode(err, errs.ServerErrorCode)
	}

	return nil
}
