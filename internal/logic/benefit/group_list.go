package benefit

import (
	"context"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupList struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupList(ctx context.Context, svcCtx *svc.ServiceContext) *GroupList {
	return &GroupList{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupList) GroupList() (resp []types.GroupListResp, err error) {
	data, err := l.svcCtx.BenefitGroupModel.FindGroupAll(l.ctx)
	if err != nil {
		return nil, err
	}

	for _, v := range data {
		resp = append(resp, types.GroupListResp{
			Id:           v.Id,
			Name:         v.Name,
			IsDeprecated: v.Intro == "废弃",
		})
	}

	return resp, nil
}
