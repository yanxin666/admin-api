package auth

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/errs"
	"e.coding.net/zmexing/nenglitanzhen/proto/ability"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type Azure struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAzure(ctx context.Context, svcCtx *svc.ServiceContext) *Azure {
	return &Azure{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Azure) Azure(req *types.AzureAuthReq) (resp *types.AzureAuthResp, err error) {
	// rpc 请求
	data, err := l.svcCtx.AbilityRPC.AuthClient.Azure(l.ctx, &ability.AzureAuthReq{
		Type: ability.AuthType(req.Type),
	})
	if err != nil {
		return nil, errs.WithMsg(err, errs.ErrRetry, "获取微软临时授权失败")
	}
	return &types.AzureAuthResp{
		Token:  data.GetToken(),
		Region: data.GetRegion(),
	}, nil
}
