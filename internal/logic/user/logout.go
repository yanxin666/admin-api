package user

import (
	"context"
	"muse-admin/internal/define"
	"muse-admin/internal/svc"
	"muse-admin/internal/tools"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout() error {
	userId := strconv.FormatInt(tools.GetUserIdByCtx(l.ctx), 10)
	_, _ = l.svcCtx.Redis.Del(define.SysPermMenuCachePrefix + userId)
	_, _ = l.svcCtx.Redis.Del(define.SysOnlineUserCachePrefix + userId)
	_, _ = l.svcCtx.Redis.Del(define.SysUserIdCachePrefix + userId)

	return nil
}
