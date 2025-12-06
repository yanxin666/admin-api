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

type GetGoodsList struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGoodsList(ctx context.Context, svcCtx *svc.ServiceContext) *GetGoodsList {
	return &GetGoodsList{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGoodsList) GetGoodsList(req *types.GoodsListReq) (resp *types.GoodsListResp, err error) {
	data, err := l.svcCtx.AbilityRPC.GoodsClient.GetGoodsList(l.ctx, &ability.SkuInfoReq{
		SkuId: req.SkuId,
	})
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	if data == nil || data.List == nil || len(data.List) == 0 {
		return nil, nil
	}

	var list []*types.GoodsInfo
	for _, item := range data.List {
		list = append(list, &types.GoodsInfo{
			Id:            item.Id,
			Name:          item.Name,
			Stock:         item.Stock,
			Description:   item.Description,
			OriginalPrice: item.OriginalPrice,
			Price:         item.Price,
			Type:          item.Type,
			Pcs:           item.Pcs,
		})
	}

	return &types.GoodsListResp{
		List: list,
	}, nil

}
