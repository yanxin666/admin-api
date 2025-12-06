package live

import (
	"context"
	"muse-admin/internal/svc"
	ctxt "muse-admin/internal/tools"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type Inspect struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInspect(ctx context.Context, svcCtx *svc.ServiceContext) *Inspect {
	return &Inspect{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Inspect) Inspect(req *types.LiveBetaUserInspectReq) (resp *types.LiveBetaUserInspectResp, err error) {
	userId := ctxt.GetUserIdByCtx(l.ctx)

	err = l.svcCtx.LiveBetaRecordModel.UpdateOperateStatus(l.ctx, req.Id, req.Status, userId)
	if err != nil {
		return nil, err
	}

	// 审批通过需要发送验证码
	if req.Status == 1 {
		// todo 这里需要一个验证码模版 @运营同学
		// err = l.svcCtx.MQProducer.SendSync(l.ctx, define.HubConsumerToPro, f)
		// if err != nil {
		// 	return nil, errs.WithMsg(err, errs.ErrPushMQ, fmt.Sprintf("发送MQ消息体失败"))
		// }
	}

	return &types.LiveBetaUserInspectResp{Result: true}, nil
}
