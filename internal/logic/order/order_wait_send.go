package order

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/proto/ability"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderWaitSend struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderWaitSend(ctx context.Context, svcCtx *svc.ServiceContext) *OrderWaitSend {
	return &OrderWaitSend{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderWaitSend) OrderWaitSend(req *types.OrderWaitSendReq) (resp *types.OrderWaitSendResp, err error) {
	param := &ability.OrderWaitSendReq{
		OrderNo:     req.OrderNo,
		SubOrderNo:  req.SubOrderNo,
		GoodsName:   req.GoodsName,
		UserName:    req.UserName,
		UserPhone:   req.UserPhone,
		GrantStatus: req.GrantStatus,
	}
	if req.Page > 0 && req.Limit > 0 {
		param.Page = &ability.Page{
			Page: req.Page,
			Size: req.Limit,
		}
	}
	timeFilter := &ability.TimeFilter{}
	if req.StartTime != "" {
		timeFilter.StartTime = req.StartTime
	}
	if req.EndTime != "" {
		timeFilter.EndTime = req.EndTime
	}
	param.TimeFilter = timeFilter

	data, err := l.svcCtx.AbilityRPC.OrderClient.OrderWaitSend(l.ctx, param)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}

	var list []*types.OrderWaitSend
	for _, item := range data.List {
		list = append(list, &types.OrderWaitSend{
			OrderNo:     item.OrderNo,
			SubOrderNo:  item.SubOrderNo,
			GoodsName:   item.GoodsName,
			Pcs:         item.Pcs,
			UserName:    item.UserName,
			UserPhone:   item.UserPhone,
			Province:    item.Province,
			City:        item.City,
			Area:        item.Area,
			Address:     item.Address,
			GrantStatus: item.GrantStatus,
			CreatedAt:   item.CreatedAt,
			PayChannel:  item.PayChannel,
		})
	}

	return &types.OrderWaitSendResp{
		Pagination: types.Pagination{
			Total: data.Total,
			Page:  req.Page,
			Limit: req.Limit,
		},
		List: list,
	}, nil
}
