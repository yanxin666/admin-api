package super_train

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"muse-admin/internal/consumer/hub/builder"
	"muse-admin/internal/consumer/hub/super_train/render"
	"muse-admin/internal/define"
	"muse-admin/internal/define/mqdef"
	"muse-admin/internal/model/hub"
	"muse-admin/internal/svc"
	"muse-admin/internal/svc/public"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"
	"regexp"
)

type ProviderSuperTrain struct {
	svcCtx *svc.ServiceContext
	f      *types.SuperTrainData
}

func NewProviderSuperTrain(svcCtx *svc.ServiceContext) builder.IBuilder {
	return &ProviderSuperTrain{
		svcCtx: svcCtx,
		f:      &types.SuperTrainData{},
	}
}

func (p *ProviderSuperTrain) ParseJson(ctx context.Context, str string) error {
	err := json.Unmarshal([]byte(str), p.f)
	if err != nil {
		logc.Errorf(ctx, "来源: 超能训练；数据解析错误；错误信息为：%v", err)
		return err
	}
	return nil
}

func (p *ProviderSuperTrain) createSnapshot(ctx context.Context) (int64, error) {
	r, _ := json.Marshal(p.f)
	// 查询版本号
	version, err := p.svcCtx.HubSuperTrainModel.FindVersionByNo(ctx, p.f.CourseNo)
	if err != nil {
		return 0, err
	}

	result, _ := p.svcCtx.HubSuperTrainModel.Insert(ctx, &hub.SuperTrainSnapshot{
		No:         p.f.CourseNo,
		Name:       p.f.CourseName,
		Status:     p.f.Status,
		Version:    version + 1,
		Data:       string(r),
		AppVersion: "",
	})

	// 获取新增ID
	aid, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	if aid == 0 {
		return 0, errors.New("新增失败，未生成自增ID")
	}

	return aid, nil
}

func isValidFormat(s string) bool {
	re := regexp.MustCompile(`^C\d{2}-\d{2}-\d{3}$`)
	return re.MatchString(s)
}

// ThirdToPro 线上环境数据处理
func (p *ProviderSuperTrain) ThirdToPro(ctx context.Context) (err error) {
	if !isValidFormat(p.f.CourseNo) {
		return errors.New("课程编号格式异常")
	}

	// 过滤测试数据
	if p.f.CourseNo == "C01-02-001" || p.f.CourseNo == "C02-03-001" || p.f.CourseNo == "C02-03-002" || p.f.CourseNo == "C02-03-003" {
		return errors.New("来源编号线上环境已存在，请勿重复")
	}

	// 数据快照
	aid, err := p.createSnapshot(ctx)
	if err != nil {
		return err
	}

	// 若数据组是测试环境过来推送的，那我们直接对接业务的测试库
	if p.svcCtx.Config.Mode == "test" {
		if err = render.NewRender(p.svcCtx).Render(ctx, p.f); err != nil {
			return err
		}

		// 推送卡片消息
		_ = public.NewBizPushMsgService(p.svcCtx).PushMsgCard(ctx, public.PushMsgData{
			Source:    define.HubSourceMap[p.f.Source],
			Mode:      define.ModeNameMap[p.svcCtx.Config.Mode],
			Operation: define.LiveOperateMap[p.f.EnvType],
			Title:     p.f.CourseName,
			Number:    p.f.CourseNo,
			Account:   p.f.Account,
		})
		return
	}

	// 由线上环境向测试环境发送MQ，直接更改测试数据
	_, err = p.svcCtx.MQProducerToTest.SendSync(ctx, mqdef.TopicHubConsumerToTest, p.f)
	if err != nil {
		return err
	}

	// 由线上环境向预发环境发送MQ，直接更改预发数据
	_, err = p.svcCtx.MQProducerToPre.SendSync(ctx, mqdef.TopicHubConsumerToPre, p.f)
	if err != nil {
		return err
	}

	// 线上审核开关是否打开
	dict, err := p.svcCtx.SysDictionaryModel.FindOneByUniqueKey(ctx, "prod_hub_super_train")
	if err != nil && !errors.Is(err, sqlc.ErrNotFound) {
		return errs.WithCode(err, errs.ServerErrorCode)
	}

	// 配置为空 || 状态为下线 || 值为不跳过
	// 线上不消费此记录，由快照表手动审核完成后执行
	if dict == nil || dict.Status == 0 || dict.Value != "1" {
		return nil
	}

	// 跳过审核
	if dict.Value == "1" {
		_ = p.svcCtx.HubSuperTrainModel.UpdateOperateStatus(ctx, aid, 1, 2, "跳过审核")
		_, err = p.svcCtx.MQProducer.SendSync(ctx, mqdef.TopicHubConsumerToPro, p.f)
		if err != nil {
			return errs.WithMsg(err, errs.ErrPushMQ, fmt.Sprintf("发送MQ消息体失败"))
		}
	}

	return nil
}

// ProToTest 直接测试环境数据处理
func (p *ProviderSuperTrain) ProToTest(ctx context.Context) error {
	return render.NewRender(p.svcCtx).Render(ctx, p.f)
}

// ProToPre 直接预发环境数据处理
func (p *ProviderSuperTrain) ProToPre(ctx context.Context) error {
	return render.NewRender(p.svcCtx).Render(ctx, p.f)
}

// ProToTestPass 测试环境无误，准备消费线上数据
func (p *ProviderSuperTrain) ProToTestPass(ctx context.Context) error {
	return render.NewRender(p.svcCtx).Render(ctx, p.f)
}
