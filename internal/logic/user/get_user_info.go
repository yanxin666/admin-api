package user

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"muse-admin/internal/define"
	"muse-admin/internal/svc"
	"muse-admin/internal/tools"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoLogic) GetUserInfo() (resp *types.UserInfoResp, err error) {
	userId := tools.GetUserIdByCtx(l.ctx)
	user, err := l.svcCtx.SysUserModel.FindOne(l.ctx, userId)
	if err != nil {
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}

	dictionary, err := l.svcCtx.SysDictionaryModel.FindOneByUniqueKey(l.ctx, "sys_pwd")
	var password string
	if dictionary.Status == define.SysEnable {
		password = dictionary.Value
	} else {
		password = define.SysNewUserDefaultPassword
	}

	return &types.UserInfoResp{
		Username:       user.Username,
		Avatar:         user.Avatar,
		IsInitPassport: user.Password == util.GenerateMD5Str(password+l.svcCtx.Config.Salt),
	}, nil
}
