package user

import (
	"context"
	"muse-admin/internal/define"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSysUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteSysUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSysUserLogic {
	return &DeleteSysUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSysUserLogic) DeleteSysUser(req *types.DeleteSysUserReq) error {
	if req.Id == define.SysSuperUserId {
		return errs.NewCode(errs.ForbiddenErrorCode)
	}

	err := l.svcCtx.SysUserModel.Delete(l.ctx, req.Id)
	if err != nil {
		return errs.WithCode(err, errs.ServerErrorCode)
	}

	return nil
}
