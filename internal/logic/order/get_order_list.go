package order

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/logz"
	"e.coding.net/zmexing/nenglitanzhen/proto/ability"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"muse-admin/internal/model/user"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderList struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOrderList(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderList {
	return &GetOrderList{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrderList) GetOrderList(req *types.GetOrderListReq) (resp *types.GetOrderListResp, err error) {
	param := &ability.OrderResearchReq{

		UserName:        req.UserName,
		OrderNo:         req.OrderNo,
		Phone:           req.Phone,
		SkuType:         req.SkuType,
		OrderStatus:     req.OrderStatus,
		PayStatus:       req.PayStatus,
		PayChannel:      req.PayChannel,
		GrantStatus:     req.GrantStatus,
		SubOrderNo:      req.SubOrderNo,
		SkuName:         req.SkuName,
		IsNeedLogistics: req.IsNeedLogistics,
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

	data, err := l.svcCtx.AbilityRPC.OrderClient.OrderResearch(l.ctx, param)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	if len(data.List) == 0 {
		return nil, nil
	}
	// 遍历userIds
	var userIds []any
	for _, v := range data.List {
		userIds = append(userIds, v.UserId)
	}
	//
	userInfos, err := user.NewUserModel(l.svcCtx.MysqlConnCenter).BatchByUserIds(l.ctx, userIds)
	if err != nil {
		logz.Errorf(l.ctx, "获取用户信息失败 err: %v", err)
		return nil, err
	}
	// 数据组装
	var list []*types.OrderResearchInfo

	for _, info := range data.List {
		researchInfo := &types.OrderResearchInfo{
			OrderNo:         info.OrderNo,
			SubOrderNo:      info.SubOrderNo,
			UserName:        info.UserName,
			Phone:           info.Phone,
			SkuName:         info.SkuName,
			SkuType:         info.SkuType,
			SkuId:           info.SkuId,
			SkuNum:          info.SkuNum,
			SkuPrice:        info.SkuPrice,
			SkuStone:        info.SkuStone,
			OrderStatus:     info.OrderStatus,
			PayStatus:       info.PayStatus,
			GrantStatus:     info.GrantStatus,
			IsNeedLogistics: info.IsNeedLogistics,
			CreatedAt:       info.CreatedAt,
		}
		if userInfo, ok := userInfos[info.UserId]; ok {
			researchInfo.UserName = userInfo.RealName
			researchInfo.Phone = userInfo.Phone
		}
		list = append(list, researchInfo)
	}

	return &types.GetOrderListResp{
		Pagination: types.Pagination{
			Total: data.Total,
			Page:  req.Page,
			Limit: req.Limit,
		},
		List: list,
	}, nil
}
