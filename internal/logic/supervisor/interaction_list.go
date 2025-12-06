package supervisor

import (
	"context"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InteractionList struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInteractionList(ctx context.Context, svcCtx *svc.ServiceContext) *InteractionList {
	return &InteractionList{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InteractionList) InteractionList(req *types.InteractionListReq) (resp *types.InteractionListResp, err error) {
	condition, err := l.svcCtx.InteractionModel.FindAllByCondition(l.ctx, req.Name)
	if err != nil {
		return nil, err
	}

	resp = &types.InteractionListResp{List: make([]types.InteractionInfo, 0)}
	for _, item := range condition {
		resp.List = append(resp.List, types.InteractionInfo{
			Id:           item.Id,
			Name:         item.Name,
			Description:  item.Description,
			TeachingData: item.TeachingData.String,
			Data:         item.Data.String,
		})
	}

	return
}

// 校验参数
func (l *InteractionList) checkParams(req *types.InteractionListReq) error {
	// todo: add your logic check params and delete this line
	// 若校验的参数过少，可删除此方法，直接在上方 InteractionList 中编写

	return nil
}
