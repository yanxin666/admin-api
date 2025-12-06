package benefit

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"muse-admin/internal/model/benefit"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"
)

type UserPhoneModify struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserPhoneModify(ctx context.Context, svcCtx *svc.ServiceContext) *UserPhoneModify {
	return &UserPhoneModify{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserPhoneModify) UserPhoneModify(req *types.UserPhoneModifyReq) (resp *types.Result, err error) {
	data, err := l.svcCtx.UserBenefitModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, errs.NewMsg(errs.ServerErrorCode, "权益列表有误").ShowMsg()
	}

	// 需要更新权益的用户所属角色有误
	if data.UserId != req.OldRoleId {
		return nil, errs.NewMsg(errs.ServerErrorCode, "当前权益与所选当前角色不一致").ShowMsg()
	}

	// 将新用户ID挂在旧用户上
	_, _ = l.svcCtx.UserBenefitModel.UpdateFillFieldsById(l.ctx, data.Id, &benefit.UserBenefitsDetail{
		UserId: req.NewRoleId,
	})

	// 更新记录表
	_, _ = l.svcCtx.UserBenefitRecordModel.UpdateFillFieldsByUserIdAndResourceId(l.ctx, req.OldRoleId, data.BenefitsResourceId, &benefit.BenefitsUserRecord{
		UserId: req.NewRoleId,
	})

	return &types.Result{
		Result: true,
	}, nil
}
