package benefit

import (
	"context"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BenefitList struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBenefitList(ctx context.Context, svcCtx *svc.ServiceContext) *BenefitList {
	return &BenefitList{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BenefitList) BenefitList() (resp []types.BenefitListResp, err error) {
	data, err := l.svcCtx.BenefitResourceModel.FindResourceAll(l.ctx)
	if err != nil {
		return nil, err
	}

	for _, v := range data {
		resp = append(resp, types.BenefitListResp{
			Id:           v.Id,
			Name:         v.Name,
			IsDeprecated: v.Description.String == "废弃",
		})
	}

	return resp, nil
}
