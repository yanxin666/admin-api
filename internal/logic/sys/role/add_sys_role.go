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

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddSysRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddSysRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddSysRoleLogic {
	return &AddSysRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddSysRoleLogic) AddSysRole(req *types.AddSysRoleReq) error {
	_, err := l.svcCtx.SysRoleModel.FindOneByUniqueKey(l.ctx, req.UniqueKey)
	if errors.Is(err, workbench.ErrNotFound) {
		if req.ParentId != define.SysTopParentId {
			_, err := l.svcCtx.SysRoleModel.FindOne(l.ctx, req.ParentId)
			if err != nil {
				return errs.NewCode(errs.ParentRoleIdErrorCode)
			}
		}

		var sysRole = new(workbench.Role)
		err = copier.Copy(sysRole, req)
		if err != nil {
			return errs.WithCode(err, errs.ServerErrorCode)
		}

		bytes, err := json.Marshal(req.PermMenuIds)
		if err != nil {
			return errs.WithCode(err, errs.ServerErrorCode)
		}

		sysRole.PermMenuIds = string(bytes)
		_, err = l.svcCtx.SysRoleModel.Insert(l.ctx, sysRole)
		if err != nil {
			return errs.WithCode(err, errs.ServerErrorCode)
		}

		return nil
	} else {

		return errs.NewCode(errs.AddRoleErrorCode)
	}
}
