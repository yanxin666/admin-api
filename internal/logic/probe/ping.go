package probe

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"muse-admin/internal/svc"
)

type Ping struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPing(ctx context.Context, svcCtx *svc.ServiceContext) *Ping {
	return &Ping{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Ping) Ping() error {
	// todo: add your logic here and delete this line

	return nil
}
