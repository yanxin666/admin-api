package lesson

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logc"
	"muse-admin/internal/define"
	"muse-admin/internal/define/mqdef"
	"muse-admin/internal/svc"
	ctxt "muse-admin/internal/tools"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"
	"time"

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

func (l *BatchInspect) BatchInspect(req *types.HubLessonBatchInspectReq) (resp *types.Result, err error) {
	userId := ctxt.GetUserIdByCtx(l.ctx)

	switch req.Source {
	case 1:
		err = l.lesson(userId, req.StartId, req.EndId)
	case 2:
		err = l.writePPT(userId, req.StartId, req.EndId)
	default:
		err = errs.NewMsg(errs.ServerErrorCode, "请填写正确来源").ShowMsg()
	}
	if err != nil {
		return nil, err
	}

	return &types.Result{Result: true}, nil
}

func (l *BatchInspect) lesson(userId, startId, endId int64) (err error) {
	data, err := l.svcCtx.HubLessonModel.FindAllByBothId(l.ctx, startId, endId)
	if err != nil {
		return err
	}

	// 多版本数据中最新的一条数据集合
	inspectMap, err := l.svcCtx.HubLessonModel.FindMaxVersion(l.ctx)
	if err != nil {
		return errs.WithCode(err, errs.ServerErrorCode)
	}

	for _, v := range data {
		// 检查 lessonNo 是否存在于 inspectMap 中，若存在，但值不相等，代表多版本的旧数据，需要过滤
		if targetValue, exists := inspectMap[cast.ToString(v.LessonNo)]; exists && v.Version != targetValue {
			continue
		}

		// 对data做一个解析
		f := &types.ScheduleData{}
		err = json.Unmarshal([]byte(v.Data), f)
		if err != nil {
			logc.Error(l.ctx, "数据解析错误为:%s 错误信息为:%v", v.Data, err)
			return err
		}

		fmt.Println(fmt.Sprintf("当前ID：%d，发送时间：%s", v.Id, time.Now().Format("2006-01-02 15:04:05")))

		err = l.svcCtx.HubLessonModel.UpdateOperateStatus(l.ctx, v.Id, 1, userId, "审核通过")
		if err != nil {
			return err
		}

		_, err = l.svcCtx.MQProducer.SendSync(l.ctx, mqdef.TopicHubConsumerToPro, f)
		if err != nil {
			return errs.WithMsg(err, errs.ErrPushMQ, fmt.Sprintf("发送MQ消息体失败"))
		}

		switch v.NodeType {
		// 小语文数据量大，设置为10秒一条
		case define.NodeType.SmallLanguage:
			time.Sleep(10 * time.Second)
		// 大语文数据量小，但是需要等待小语文入库，设置为2秒一条
		case define.NodeType.BigLanguage:
			time.Sleep(2 * time.Second)
		// 默认设置为500毫秒一条
		default:
			time.Sleep(500 * time.Millisecond) // 设置传输间隔为500毫秒，1秒内发送2条
		}

		fmt.Println(fmt.Sprintf("当前ID：%d，处理完成：%s", v.Id, time.Now().Format("2006-01-02 15:04:05")))
	}

	return nil
}

func (l *BatchInspect) writePPT(userId, startId, endId int64) (err error) {
	data, err := l.svcCtx.HubWritePPTModel.FindAllByBothId(l.ctx, startId, endId)
	if err != nil {
		return err
	}

	for _, v := range data {
		f := &types.WritePPTData{}
		err = json.Unmarshal([]byte(v.Data), f)
		if err != nil {
			logc.Error(l.ctx, "数据解析错误为:%s 错误信息为:%v", v.Data, err)
			return err
		}

		err = l.svcCtx.HubWritePPTModel.UpdateOperateStatus(l.ctx, v.Id, 1, userId, "审核通过")
		if err != nil {
			return err
		}

		_, err = l.svcCtx.MQProducer.SendSync(l.ctx, mqdef.TopicHubConsumerToPro, f)
		if err != nil {
			return errs.WithMsg(err, errs.ErrPushMQ, fmt.Sprintf("发送MQ消息体失败"))
		}

		time.Sleep(100 * time.Millisecond) // 设置传输间隔为100毫秒，1秒内发送10条，1分内发送600个
	}

	return nil
}
