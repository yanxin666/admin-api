package benefit

import (
	"context"
	"muse-admin/internal/define"
	"muse-admin/internal/model/benefit"
	"muse-admin/pkg/errs"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type Reissue struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReissue(ctx context.Context, svcCtx *svc.ServiceContext) *Reissue {
	return &Reissue{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Reissue) Reissue(req *types.ReissueReq) (resp *types.Result, err error) {
	data, err := l.svcCtx.UserBenefitModel.FindOne(l.ctx, req.UserBenefitId)
	if err != nil {
		return nil, errs.NewMsg(errs.ServerErrorCode, "权益列表有误").ShowMsg()
	}

	// 权益生效中，不能重新发放
	if data.Status == define.BenefitsStatus.Effective {
		return nil, errs.NewMsg(errs.ServerErrorCode, "当前权益正在生效中，无需重新发放").ShowMsg()
	}

	// 将权益状态改为生效中
	_, _ = l.svcCtx.UserBenefitModel.UpdateFillFieldsById(l.ctx, data.Id, &benefit.UserBenefitsDetail{
		Status: define.BenefitsStatus.Effective,
	})

	return &types.Result{
		Result: true,
	}, nil
}
