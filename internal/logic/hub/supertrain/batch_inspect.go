package supertrain

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/logz"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"muse-admin/internal/define/mqdef"
	ctxt "muse-admin/internal/tools"
	"muse-admin/pkg/errs"
	"time"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchInspect struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchInspect(ctx context.Context, svcCtx *svc.ServiceContext) *BatchInspect {
	return &BatchInspect{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchInspect) BatchInspect(req *types.HubSuperTrainBatchInspectReq) (resp *types.Result, err error) {

	userId := ctxt.GetUserIdByCtx(l.ctx)

	data, err := l.svcCtx.HubSuperTrainModel.FindAllByBothId(l.ctx, req.StartId, req.EndId)
	if err != nil {
		return nil, err
	}

	// 多版本数据中最新的一条数据集合
	inspectMap, err := l.svcCtx.HubSuperTrainModel.FindMaxVersion(l.ctx)
	if err != nil {
		return nil, err
	}

	for _, v := range data {
		// 检查 lessonNo 是否存在于 inspectMap 中，若存在，但值不相等，代表多版本的旧数据，需要过滤
		if targetValue, exists := inspectMap[v.No]; exists && v.Version != targetValue {
			continue
		}

		// 对data做一个解析
		f := &types.ScheduleData{}
		err = json.Unmarshal([]byte(v.Data), f)
		if err != nil {
			logc.Error(l.ctx, "数据解析错误为:%s 错误信息为:%v", v.Data, err)
			return nil, err
		}

		logz.Infof(l.ctx, "当前ID：%d，发送时间：%s", v.Id, time.Now().Format("2006-01-02 15:04:05"))

		err = l.svcCtx.HubSuperTrainModel.UpdateOperateStatus(l.ctx, v.Id, 1, userId, "审核通过")
		if err != nil {
			return nil, err
		}

		_, err = l.svcCtx.MQProducer.SendSync(l.ctx, mqdef.TopicHubConsumerToPro, f)
		if err != nil {
			return nil, errs.WithMsg(err, errs.ErrPushMQ, fmt.Sprintf("发送MQ消息体失败"))
		}

		// 默认设置为500毫秒一条
		time.Sleep(500 * time.Millisecond) // 设置传输间隔为500毫秒，1秒内发送2条

		logz.Infof(l.ctx, "当前ID：%d，处理完成：%s", v.Id, time.Now().Format("2006-01-02 15:04:05"))
	}

	return &types.Result{Result: true}, nil
}
