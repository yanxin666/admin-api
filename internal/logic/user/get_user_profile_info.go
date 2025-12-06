package user

import (
	"context"
	"muse-admin/internal/svc"
	"muse-admin/internal/tools"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserProfileInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserProfileInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserProfileInfoLogic {
	return &GetUserProfileInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserProfileInfoLogic) GetUserProfileInfo() (resp *types.UserProfileInfoResp, err error) {
	userId := tools.GetUserIdByCtx(l.ctx)
	user, err := l.svcCtx.SysUserModel.FindOne(l.ctx, userId)
	if err != nil {
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}

	return &types.UserProfileInfoResp{
		Username: user.Username,
		Nickname: user.Nickname,
		Gender:   user.Gender,
		Email:    user.Email,
		Mobile:   user.Mobile,
		Remark:   user.Remark,
		Avatar:   user.Avatar,
	}, nil
}
