package user

import (
	"context"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"muse-admin/pkg/other"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGenerateAvatarLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGenerateAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGenerateAvatarLogic {
	return &GetGenerateAvatarLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGenerateAvatarLogic) GetGenerateAvatar() (resp *types.GenerateAvatarResp, err error) {
	return &types.GenerateAvatarResp{
		AvatarUrl: other.AvatarUrl(),
	}, nil
}
