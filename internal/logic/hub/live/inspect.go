package live

import (
	"context"
	"muse-admin/internal/svc"
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

func (l *Inspect) Inspect(req *types.HubLiveInspectReq) (resp *types.Result, err error) {
	// userId := ctxt.GetUserIdByCtx(l.ctx)
	//
	// data, err := l.svcCtx.HubLessonModel.FindOne(l.ctx, req.LiveNo)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// // 对data做一个解析
	// f := &types.ScheduleData{}
	// err = json.Unmarshal([]byte(data.Data), f)
	// if err != nil {
	// 	logc.Error(l.ctx, "数据解析错误为:%s 错误信息为:%v", data.Data, err)
	// 	return nil, err
	// }
	//
	// // data.OperateStatus = 1
	// // data.OperateId = userId
	//
	// err = l.svcCtx.HubLessonModel.UpdateOperateStatus(l.ctx, data.Id, 1, userId, "审核通过")
	// if err != nil {
	// 	return nil, err
	// }
	//
	// err = l.svcCtx.MQProducer.SendSync(l.ctx, define.HubConsumerToPro, f)
	// if err != nil {
	// 	return nil, errs.WithMsg(err, errs.ErrPushMQ, fmt.Sprintf("发送MQ消息体失败"))
	// }

	return &types.Result{Result: true}, nil
}
