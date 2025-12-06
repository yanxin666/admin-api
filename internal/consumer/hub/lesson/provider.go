package lesson

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"muse-admin/internal/consumer/hub/builder"
	"muse-admin/internal/consumer/hub/lesson/render"
	"muse-admin/internal/define/mqdef"
	"muse-admin/internal/model/hub"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"
)

type ProviderLesson struct {
	svcCtx *svc.ServiceContext
	f      *types.ScheduleData
}

func NewProviderLesson(svcCtx *svc.ServiceContext) builder.IBuilder {
	return &ProviderLesson{
		svcCtx: svcCtx,
		f:      &types.ScheduleData{},
	}
}

func (p *ProviderLesson) ParseJson(ctx context.Context, str string) error {
	err := json.Unmarshal([]byte(str), p.f)
	if err != nil {
		logc.Errorf(ctx, "来源：写作课程；数据解析错误；错误信息为：%v", err)
		return err
	}
	return nil
}

func (p *ProviderLesson) createSnapshot(ctx context.Context) (int64, error) {
	r, _ := json.Marshal(p.f)
	// 查询版本号
	version, err := p.svcCtx.HubLessonModel.FindVersionByLessonNo(ctx, p.f.Id)
	if err != nil {
		return 0, err
	}

	result, _ := p.svcCtx.HubLessonModel.Insert(ctx, &hub.LessonSnapshot{
		Id:            0,
		LessonNo:      p.f.Id,
		NodeType:      p.f.NodeType,
		Name:          p.f.CourseName,
		PointName:     p.f.PointName,
		Status:        p.f.Status,
		Grade:         p.f.Grade,
		Version:       version + 1,
		OperateStatus: 0,
		OperateId:     0,
		Data:          string(r),
		AppVersion:    "",
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

// ThirdToPro 线上环境数据处理
func (p *ProviderLesson) ThirdToPro(ctx context.Context) (err error) {
	if p.f.Id == 0 {
		return errors.New("来源编号异常")
	}

	// 预插入头部的例题
	render.AddExampleQuestion(p.f)

	// 数据快照
	aid, err := p.createSnapshot(ctx)
	if err != nil {
		return err
	}

	// 由线上环境向测试环境发送MQ，直接更改测试数据
	_, err = p.svcCtx.MQProducerToTest.SendSync(ctx, mqdef.TopicHubConsumerToTest, p.f)
	if err != nil {
		return err
	}

	// 线上审核开关是否打开
	dict, err := p.svcCtx.SysDictionaryModel.FindOneByUniqueKey(ctx, "prod_hub_lesson")
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
		_ = p.svcCtx.HubLessonModel.UpdateOperateStatus(ctx, aid, 1, 2, "跳过审核")
		_, err = p.svcCtx.MQProducer.SendSync(ctx, mqdef.TopicHubConsumerToPro, p.f)
		if err != nil {
			return errs.WithMsg(err, errs.ErrPushMQ, fmt.Sprintf("发送MQ消息体失败"))
		}
	}
	return nil
}

// ProToTest 直接测试环境数据处理
func (p *ProviderLesson) ProToTest(ctx context.Context) error {
	return render.NewSchedule(p.svcCtx).Render(ctx, p.f)
}

// ProToPre 直接预发环境数据处理
func (p *ProviderLesson) ProToPre(ctx context.Context) error {
	return render.NewSchedule(p.svcCtx).Render(ctx, p.f)
}

// ProToTestPass 测试环境无误，准备消费线上数据
func (p *ProviderLesson) ProToTestPass(ctx context.Context) error {
	return render.NewSchedule(p.svcCtx).Render(ctx, p.f)
}
