package hub

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/mq/types"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/logz"
	"encoding/json"
	"errors"
	"fmt"
	rmq "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/zeromicro/go-zero/core/logc"
	"muse-admin/internal/consumer/hub/builder"
	"muse-admin/internal/consumer/hub/lesson"
	"muse-admin/internal/consumer/hub/live"
	"muse-admin/internal/consumer/hub/super_train"
	"muse-admin/internal/consumer/hub/write_ppt"
	"muse-admin/internal/define"
	"muse-admin/internal/svc"
)

type DataHub struct {
	svcCtx *svc.ServiceContext
	gen    *builder.Generator
	source int64
	data   string
}

func NewDataHub(svcCtx *svc.ServiceContext, gen *builder.Generator) types.Consumer {
	return &DataHub{
		svcCtx: svcCtx,
		gen:    gen,
	}
}

func (s *DataHub) Explain(ctx context.Context, data interface{}) error {
	logc.Infof(ctx, "[DataHub_Explain] 来源：%s", s.gen.GetReqFrom())

	jsonString, ok := data.(string)
	if !ok {
		logc.Infof(ctx, "data is not a string")
	}

	// 获取source来源
	source, err := getSourceValue(jsonString)
	if err != nil {
		return errors.New(fmt.Sprintf("来源Source有误，Err:%v", err))
	}

	s.source = source
	s.data = jsonString

	return nil
}

// Execute 区分来源以及去向，角色定位于调度器
func (s *DataHub) Execute(ctx context.Context) (err error) {
	mapHub := map[int64]builder.IBuilder{
		define.Lesson:     lesson.NewProviderLesson(s.svcCtx),
		define.WritePPT:   write_ppt.NewProviderWritePPT(s.svcCtx),
		define.AILive:     live.NewProviderLive(s.svcCtx),
		define.SuperTrain: super_train.NewProviderSuperTrain(s.svcCtx),
	}

	var title string
	defer func() {
		if err != nil {
			logz.Errorf(ctx, "数据有误，err=%v data=%v \n", err, s.data)
			_ = define.AlarmMqCustomer(ctx, s.svcCtx.Config.Mode, define.Alarm.Wpf, s.svcCtx.Config.Cls.TopicID, title, s.gen.GetReqFrom(), s.data, err)
			// 报警后消费此条mq
			err = nil
			return
		}
	}()

	// builder构建 && 执行
	iBuilder, ok := mapHub[s.source]
	if !ok {
		title = "MQ获取执行对象失败 -> 来源Source不存在"
		return fmt.Errorf("source: %d 不存在", s.source)
	}
	err = iBuilder.ParseJson(ctx, s.data)
	if err != nil {
		return err
	}

	renderMap := map[string]func(ctx context.Context) error{
		"ThirdToOnline":  iBuilder.ThirdToPro,
		"OnlineToTest":   iBuilder.ProToTest,
		"OnlineToPre":    iBuilder.ProToPre,
		"ManageToOnline": iBuilder.ProToTestPass,
	}
	// 获取请求来源并执行相应的操作
	action, exists := renderMap[s.gen.GetReqFrom()]
	if !exists {
		title = "MQ开始消费失败 -> 未注册"
		return fmt.Errorf("source: %s 可执行的renderMap不存在", s.gen.GetReqFrom())
	}

	// 执行动作并处理错误
	err = action(ctx)
	if err != nil {
		title = "MQ执行失败 -> 内部处理逻辑有误"
		return err
	}

	// 将此条mq消费掉
	return nil
}

func (s *DataHub) ExecuteBatch(ctx context.Context) ([]*rmq.MessageView, error) {
	return nil, nil
}

// getSourceValue 判断 data 中的 Source 是否存在，并返回其值
func getSourceValue(jsonString string) (int64, error) {
	var dataMap map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &dataMap)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("Error decoding JSON:%v", err))
	}
	// 获取source值
	source := dataMap["source"].(float64) // JSON解码时数字会被转换为float64

	return int64(source), nil
}
