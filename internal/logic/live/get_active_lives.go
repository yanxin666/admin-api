package live

import (
	"context"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"e.coding.net/zmexing/nenglitanzhen/biz-lib/logz"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetActiveLives struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetActiveLives(ctx context.Context, svcCtx *svc.ServiceContext) *GetActiveLives {
	return &GetActiveLives{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetActiveLives) GetActiveLives() (resp *types.GetActiveLivesResp, err error) {
	Schedule, err := l.svcCtx.SupeScheduleModel.FindActiveSchedule(l.ctx)
	if err != nil {
		logz.Errorf(l.ctx, "获取当前活动排课失败，Err:%s", err)
		return nil, err
	}
	var list []types.ActiveLiveInfo
	for _, item := range Schedule {
		list = append(list, types.ActiveLiveInfo{
			StreamId: item.StreamId,
			Name:     item.Name,
		})
	}
	return &types.GetActiveLivesResp{
		List: list,
	}, nil
}

// 校验参数
