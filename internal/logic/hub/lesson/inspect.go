package lesson

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"muse-admin/internal/define/mqdef"
	"muse-admin/internal/svc"
	ctxt "muse-admin/internal/tools"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"
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

func (l *Inspect) Inspect(req *types.HubLessonInspectReq) (resp *types.Result, err error) {
	userId := ctxt.GetUserIdByCtx(l.ctx)

	data, err := l.svcCtx.HubLessonModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, err
	}

	// 对data做一个解析
	f := &types.ScheduleData{}
	err = json.Unmarshal([]byte(data.Data), f)
	if err != nil {
		logc.Error(l.ctx, "数据解析错误为:%s 错误信息为:%v", data.Data, err)
		return nil, err
	}

	// data.OperateStatus = 1
	// data.OperateId = userId

	err = l.svcCtx.HubLessonModel.UpdateOperateStatus(l.ctx, data.Id, 1, userId, "审核通过")
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.MQProducer.SendSync(l.ctx, mqdef.TopicHubConsumerToPro, f)
	if err != nil {
		return nil, errs.WithMsg(err, errs.ErrPushMQ, fmt.Sprintf("发送MQ消息体失败"))
	}

	return &types.Result{Result: true}, nil
}
