package auth

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/errs"
	"e.coding.net/zmexing/nenglitanzhen/proto/ability"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type Tencent struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTencent(ctx context.Context, svcCtx *svc.ServiceContext) *Tencent {
	return &Tencent{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Tencent) Tencent(req *types.TencentAuthReq) (resp *types.TencentAuthResp, err error) {
	// rpc 请求
	data, err := l.svcCtx.AbilityRPC.AuthClient.Tencent(l.ctx, &ability.TencentAuthReq{
		Type: ability.AuthType(req.Type),
	})
	if err != nil {
		return nil, errs.WithMsg(err, errs.ErrRetry, "获取腾讯临时授权失败")
	}
	return &types.TencentAuthResp{
		TmpSecretID:  data.TmpSecretId,
		SessionToken: data.SessionToken,
		RequestId:    data.RequestId,
		ExpiredTime:  data.ExpiredTime,
		Bucket:       data.Bucket,
		Region:       data.Region,
		FileDir:      data.FileDir,
		StartTime:    data.StartTime,
		TmpSecretKey: data.TmpSecretKey,
		StaticUrl:    data.Host,
	}, nil
}
