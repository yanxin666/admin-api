package benefit

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"e.coding.net/zmexing/nenglitanzhen/proto/ability"
	"muse-admin/pkg/errs"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserBenefitAdd struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserBenefitAdd(ctx context.Context, svcCtx *svc.ServiceContext) *UserBenefitAdd {
	return &UserBenefitAdd{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserBenefitAdd) UserBenefitAdd(req *types.UserBenefitAddReq) (resp *types.Result, err error) {
	if err = l.checkParams(req); err != nil {
		return nil, err
	}

	_, err = l.svcCtx.AbilityRPC.BenefitClient.BenefitsGrant(l.ctx, &ability.BenefitsGrantReq{
		Phone:           req.Phone,
		BenefitsGroupId: req.GroupId,
		ChannelOrder:    req.OrderId,
		Source:          int64(100),
		Count:           req.Num,
	})
	if err != nil {
		return nil, err
	}

	return &types.Result{
		Result: true,
	}, nil
}

// 校验参数
func (l *UserBenefitAdd) checkParams(req *types.UserBenefitAddReq) error {
	if !util.CheckMobile(req.Phone) {
		return errs.NewMsg(errs.ErrCodeParamsAbnormal, "手机号格式不正确").ShowMsg()
	}

	return nil
}
