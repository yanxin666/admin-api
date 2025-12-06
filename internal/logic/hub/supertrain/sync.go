package supertrain

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/logz"
	"encoding/json"
	"fmt"
	"muse-admin/internal/consumer/hub/builder"
	"muse-admin/internal/consumer/hub/lesson"
	"muse-admin/internal/consumer/hub/live"
	"muse-admin/internal/consumer/hub/super_train"
	"muse-admin/internal/consumer/hub/write_ppt"
	"muse-admin/internal/define"
	"muse-admin/pkg/errs"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type Sync struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSync(ctx context.Context, svcCtx *svc.ServiceContext) *Sync {
	return &Sync{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Sync) Sync(req *types.HubSuperTrainSyncReq) (resp *types.Result, err error) {

	mapHub := map[int64]builder.IBuilder{
		define.Lesson:     lesson.NewProviderLesson(l.svcCtx),
		define.WritePPT:   write_ppt.NewProviderWritePPT(l.svcCtx),
		define.AILive:     live.NewProviderLive(l.svcCtx),
		define.SuperTrain: super_train.NewProviderSuperTrain(l.svcCtx),
	}

	mqData := &types.SuperTrainData{}
	err = json.Unmarshal([]byte(req.Data), mqData)

	var title string
	defer func() {
		if err != nil {
			logz.Errorf(l.ctx, "数据有误，err=%v data=%v \n", err, req.Data)
			_ = define.AlarmMqCustomer(l.ctx, l.svcCtx.Config.Mode, define.Alarm.Wpf, l.svcCtx.Config.Cls.TopicID, title, "HTTP请求", req.Data, err)
			return
		}
	}()

	// builder构建 && 执行
	iBuilder, ok := mapHub[mqData.Source]
	if !ok {
		title = "MQ获取执行对象失败 -> 来源Source不存在"
		return nil, errs.NewMsg(errs.ErrCodeParamsAbnormal, fmt.Sprintf("source: %d 不存在", mqData.Source))
	}
	err = iBuilder.ParseJson(l.ctx, req.Data)
	if err != nil {
		return nil, errs.WithCode(err, errs.ErrCodeParamsAbnormal)
	}

	// 执行动作并处理错误
	err = iBuilder.ProToTestPass(l.ctx)
	if err != nil {
		title = "MQ执行失败 -> 内部处理逻辑有误"
		return nil, errs.WithCode(err, errs.ErrCodeParamsAbnormal)
	}

	return
}
