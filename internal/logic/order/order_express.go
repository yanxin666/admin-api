package order

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/proto/ability"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderExpress struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderExpress(ctx context.Context, svcCtx *svc.ServiceContext) *OrderExpress {
	return &OrderExpress{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderExpress) OrderExpress(req *types.OrderExpressReq) (resp *types.OrderExpressResp, err error) {
	if err = l.checkParams(req); err != nil {
		return nil, err
	}

	_, err = l.svcCtx.AbilityRPC.OrderClient.OrderExpress(l.ctx, &ability.OrderExpressReq{SubOrderNo: req.SubOrderNo})
	if err != nil {
		return nil, err
	}

	return
}

// 校验参数
func (l *OrderExpress) checkParams(req *types.OrderExpressReq) error {
	// todo: add your logic check params and delete this line
	// 若校验的参数过少，可删除此方法，直接在上方 OrderExpress 中编写

	return nil
}
