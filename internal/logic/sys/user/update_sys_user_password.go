package user

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"github.com/jinzhu/copier"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"muse-admin/internal/define"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"
)

type UpdateSysUserPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSysUserPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSysUserPasswordLogic {
	return &UpdateSysUserPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSysUserPasswordLogic) UpdateSysUserPassword(req *types.UpdateSysUserPasswordReq) error {
	sysUser, err := l.svcCtx.SysUserModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return errs.NewCode(errs.UserIdErrorCode)
	}

	err = copier.Copy(sysUser, req)
	if err != nil {
		return errs.WithCode(err, errs.ServerErrorCode)
	}

	sysUser.Password = util.GenerateMD5Str(req.Password + l.svcCtx.Config.Salt)
	err = l.svcCtx.SysUserModel.Update(l.ctx, sysUser)
	if err != nil {
		return errs.WithCode(err, errs.ServerErrorCode)
	}

	_, _ = l.svcCtx.Redis.Del(define.SysPermMenuCachePrefix + cast.ToString(req.Id))
	_, _ = l.svcCtx.Redis.Del(define.SysOnlineUserCachePrefix + cast.ToString(req.Id))
	_, _ = l.svcCtx.Redis.Del(define.SysUserIdCachePrefix + cast.ToString(req.Id))

	return nil
}
