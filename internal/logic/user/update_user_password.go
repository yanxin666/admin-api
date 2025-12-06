package user

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"github.com/jefferyjob/go-easy-utils/v2/validUtil"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"muse-admin/internal/define"
	"muse-admin/internal/svc"
	"muse-admin/internal/tools"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"
)

type UpdateUserPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserPasswordLogic {
	return &UpdateUserPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserPasswordLogic) UpdateUserPassword(req *types.UpdatePasswordReq) error {
	dictionary, err := l.svcCtx.SysDictionaryModel.FindOneByUniqueKey(l.ctx, "sys_ch_pwd")
	if dictionary.Status == define.SysDisable {
		return errs.NewCode(errs.ForbiddenErrorCode)
	}

	userId := tools.GetUserIdByCtx(l.ctx)
	user, err := l.svcCtx.SysUserModel.FindOne(l.ctx, userId)
	if err != nil {
		return errs.WithCode(err, errs.ServerErrorCode)
	}

	if user.Password != util.GenerateMD5Str(req.OldPassword+l.svcCtx.Config.Salt) {
		return errs.NewCode(errs.PasswordErrorCode)
	}

	// 密码长度在6-20个字符之间，只包含数字、字母和下划线
	if !validUtil.IsPassword(req.NewPassword) {
		return errs.WithCode(err, errs.PassportFormatErrorCode)
	}

	user.Password = util.GenerateMD5Str(req.NewPassword + l.svcCtx.Config.Salt)
	err = l.svcCtx.SysUserModel.Update(l.ctx, user)
	if err != nil {
		return errs.WithCode(err, errs.ServerErrorCode)
	}

	_, _ = l.svcCtx.Redis.Del(define.SysPermMenuCachePrefix + cast.ToString(userId))
	_, _ = l.svcCtx.Redis.Del(define.SysOnlineUserCachePrefix + cast.ToString(userId))
	_, _ = l.svcCtx.Redis.Del(define.SysUserIdCachePrefix + cast.ToString(userId))

	return nil
}
