package role

import (
	"context"
	"muse-admin/internal/define"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSysRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteSysRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSysRoleLogic {
	return &DeleteSysRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSysRoleLogic) DeleteSysRole(req *types.DeleteSysRoleReq) error {
	if req.Id == define.SysSuperRoleId {
		return errs.NewCode(errs.ForbiddenErrorCode)
	}

	roleList, _ := l.svcCtx.SysRoleModel.FindSubRole(l.ctx, req.Id)
	if len(roleList) != 0 {
		return errs.NewCode(errs.DeleteRoleErrorCode)
	}

	count, _ := l.svcCtx.SysUserModel.FindCountByRoleId(l.ctx, req.Id)
	if count != 0 {
		return errs.NewCode(errs.RoleIsUsingErrorCode)
	}

	err := l.svcCtx.SysRoleModel.Delete(l.ctx, req.Id)
	if err != nil {
		return errs.WithCode(err, errs.ServerErrorCode)
	}

	return nil
}
