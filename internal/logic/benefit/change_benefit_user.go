package benefit

import (
	"context"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"

	"e.coding.net/zmexing/nenglitanzhen/biz-lib/logz"
	"e.coding.net/zmexing/nenglitanzhen/proto/passport"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ChangeBenefitUser struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangeBenefitUser(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeBenefitUser {
	return &ChangeBenefitUser{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangeBenefitUser) ChangeBenefitUser(req *types.ChangeBenefitUserReq) (resp *types.Result, err error) {
	resp = &types.Result{
		Result: true,
	}

	respRoleList, err := l.svcCtx.PassportRPC.NoAuthClient.GetRoleListByBaseUserId(l.ctx, &passport.BaseUserIdReq{
		BaseUserId: req.BaseUserId,
	})
	if err != nil || respRoleList == nil || len(respRoleList.UserData) == 0 {
		logz.Errorf(l.ctx, "ChangeBenefitUser GetRoleListByBaseUserId err: %v respRoleList:%v", err, respRoleList)
		resp.Result = false
		return resp, errs.NewMsg(errs.ServerErrorCode, "用户数据不存在").ShowMsg()
	}

	var exist bool
	for _, u := range respRoleList.UserData {
		if u.UserId == req.TargetUserId {
			exist = true
			break
		}
	}

	if !exist {
		logz.Warnf(l.ctx, "ChangeBenefitUser 目标用户不属于该主账号用户 benefitId:%d, baseUserId:%d, targetUserId:%d", req.BenefitId, req.BaseUserId, req.TargetUserId)
		resp.Result = false
		return resp, nil
	}

	// 查询权益信息
	benefit, err := l.svcCtx.UserBenefitModel.FindOne(l.ctx, req.BenefitId)
	if err != nil {
		resp.Result = false
		logz.Infof(l.ctx, "ChangeBenefitUser UserBenefitModel FindOne err: %v ", err)
		return nil, errs.NewMsg(errs.ServerErrorCode, "未找到用户权益").ShowMsg()
	}

	benefitTarget, err := l.svcCtx.UserBenefitModel.FindOneByUserIdBenefitsResourceId(l.ctx, req.TargetUserId, benefit.BenefitsResourceId)
	if err == nil && benefitTarget != nil {
		resp.Result = false
		logz.Infof(l.ctx, "ChangeBenefitUser 目标用户已存在该权益资源 benefitId:%d, targetUserId:%d BenefitsResourceId:%d", req.BenefitId, req.TargetUserId, benefit.BenefitsResourceId)
		return resp, errs.NewMsg(errs.ServerErrorCode, "目标用户已存在该权益资源").ShowMsg()
	}

	oldUserId := benefit.UserId
	record, err := l.svcCtx.UserBenefitRecordModel.FindOneByUserIdAndBenefitId(l.ctx, oldUserId, req.BenefitId)
	if err != nil {
		resp.Result = false
		logz.Infof(l.ctx, "ChangeBenefitUser UserBenefitRecordModel FindOne err: %v ", err)
		return nil, errs.NewMsg(errs.ServerErrorCode, "未找到用户权益记录").ShowMsg()
	}

	order, err := l.svcCtx.ChannelOrderModel.FindOneByOrderNo(l.ctx, record.OrderNo)
	if err != nil {
		resp.Result = false
		logz.Infof(l.ctx, "ChangeBenefitUser ChannelOrderModel FindOne err: %v ", err)
		return nil, errs.NewMsg(errs.ServerErrorCode, "未找到用户权益相关订单").ShowMsg()
	}

	logz.Infof(l.ctx, "ChangeBenefitUser benefitId:%d, baseUserId:%d, oldUserId:%d, targetUserId:%d ", benefit.Id, req.BaseUserId, oldUserId, req.TargetUserId)
	err = l.svcCtx.MysqlConnAbility.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		err := l.svcCtx.UserBenefitModel.UpdateUserIdWithTx(l.ctx, session, req.TargetUserId, benefit.Id)
		if err != nil {
			return err
		}

		err = l.svcCtx.UserBenefitRecordModel.UpdateUserIdWithTx(l.ctx, session, req.TargetUserId, oldUserId, benefit.Id)
		if err != nil {
			return err
		}

		err = l.svcCtx.ChannelSubOrderModel.UpdateUserIdWithTx(l.ctx, session, req.TargetUserId, oldUserId, order.Id, benefit.BenefitsResourceId)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		resp.Result = false
		return nil, errs.NewMsg(errs.ServerErrorCode, "更新数据库失败").ShowMsg()
	}

	return resp, nil
}
