package live

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/logz"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/third/feishu"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/trace"
	"muse-admin/internal/consumer/hub/builder"
	"muse-admin/internal/consumer/hub/live/render"
	"muse-admin/internal/define"
	"muse-admin/internal/model/hub"
	"muse-admin/internal/model/workbench"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"
)

type ProviderLive struct {
	svcCtx *svc.ServiceContext
	f      *types.LiveData
}

const (
	liveCourseLock       = "muse-admin:live:course:%s" // 推送行为
	liveCourseLockExpire = 10                          // 行为失效时间 单位(s)
)

func NewProviderLive(svcCtx *svc.ServiceContext) builder.IBuilder {
	return &ProviderLive{
		svcCtx: svcCtx,
		f:      &types.LiveData{},
	}
}

func (p *ProviderLive) ParseJson(ctx context.Context, str string) error {
	err := json.Unmarshal([]byte(str), p.f)
	if err != nil {
		logc.Errorf(ctx, "来源：AI直播课；数据解析错误；错误信息为：%v", err)
		return err
	}

	if p.f.EnvType != 1 && p.f.EnvType != 2 {
		return errors.New(fmt.Sprintf("课堂编号：%s，推送环境有误，请检查！", p.f.LiveNo))
	}
	if p.f.ModeType != 1 && p.f.ModeType != 2 {
		return errors.New(fmt.Sprintf("课堂编号：%s，推送模式有误，请检查模版模式和定制工厂模式的枚举值！", p.f.LiveNo))
	}
	if p.f.Content == "" {
		return errors.New(fmt.Sprintf("课堂编号：%s，content字段为空，请检查", p.f.LiveNo))
	}
	if len(p.f.IntervalDuration) != 2 {
		return errors.New(fmt.Sprintf("课堂编号：%s，寒暄随机间隔区间数：%#v，必须为两个种子！", p.f.LiveNo, p.f.IntervalDuration))
	}
	if !util.IsJSON(p.f.Draft) {
		return errors.New(fmt.Sprintf("课堂编号：%s，字幕格式有误，非json格式！", p.f.LiveNo))
	}
	if !util.IsJSON(p.f.BlockTime) {
		return errors.New(fmt.Sprintf("课堂编号：%s，气口格式有误，非json格式！", p.f.LiveNo))
	}
	if !util.IsJSON(p.f.Content) {
		return errors.New(fmt.Sprintf("课堂编号：%s，课中内容格式有误，非json格式！", p.f.LiveNo))
	}
	// 若为工厂模式，老师名称、音色、音频、举手等字段都需要有，若没有就报错
	if p.f.ModeType == 2 {
		if p.f.TeacherName == "" || p.f.TeacherNickname == "" || p.f.TeacherToneModel == "" || p.f.TeacherTone == "" || p.f.TeacherToneType == 0 || p.f.HandAudio == "" {
			return errors.New(fmt.Sprintf("课堂编号：%s，定制工厂的老师必填项缺失，请检查！", p.f.LiveNo))
		}
	}
	if p.f.TeacherGif == "" || p.f.TeacherPng == "" {
		return errors.New(fmt.Sprintf("课堂编号：%s，课程缺少关键信息：老师讲课的动图和静图，请检查！", p.f.LiveNo))
	}

	// // 导航栏后续考虑是否强制性效验
	// if len(p.f.Navigate) <= 0 {
	// 	return errors.New(fmt.Sprintf("课堂编号：%s，环节进度个数有误，请检查", p.f.LiveNo))
	// }

	// 分布式锁Key
	lockKey := fmt.Sprintf(liveCourseLock, p.f.LiveNo)
	// 分布式锁value，用于防止解锁其他线程的锁
	lockValue := util.GetUUIDWithoutHyphen()
	if err = p.getLock(ctx, lockKey, lockValue); err != nil {
		return err
	}

	return nil
}

func (p *ProviderLive) createSnapshot(ctx context.Context) (int64, error) {
	r, _ := json.Marshal(p.f)
	// 查询版本号
	version, err := p.svcCtx.HubLiveModel.FindVersionByLessonNo(ctx, p.f.LiveNo)
	if err != nil {
		return 0, err
	}

	result, _ := p.svcCtx.HubLiveModel.Insert(ctx, &hub.LiveSnapshot{
		LessonNo:   p.f.LiveNo,
		Name:       p.f.Name,
		Status:     p.f.ReviewStatus,
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

// ThirdToPro 线上环境数据处理
func (p *ProviderLive) ThirdToPro(ctx context.Context) (err error) {
	if p.f.LiveNo == "" {
		return errors.New("来源编号异常")
	}

	// 数据快照
	aid, err := p.createSnapshot(ctx)
	if err != nil {
		return err
	}

	fmt.Sprintln(aid)

	// 如果三方来源是测试环境过来的，那我们直接对接业务的测试库
	if p.svcCtx.Config.Mode == "test" {
		// 错误通知
		if err = render.NewRender(p.svcCtx).Render(ctx, p.f); err != nil {
			return err
		}

		// 推送消息
		userId := p.getUserId(ctx)
		if userId != "" {
			_ = p.PushMessageByCard(ctx, userId)
		} else {
			_ = p.PushMessageToAdmin(ctx)
		}
		return
	}

	// // 预览需要打到测试库供数据组预览
	// if p.f.EnvType == 1 {
	// 	// 由线上环境向测试环境发送MQ，直接更改测试数据
	// 	err = p.svcCtx.MQProducerToTest.SendSync(ctx, define.HubConsumerToTest, p.f)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// // 线上审核开关是否打开
	// dict, err := p.svcCtx.SysDictionaryModel.FindOneByUniqueKey(ctx, "prod_hub_write_ppt")
	// if err != nil && !errors.Is(err, sqlc.ErrNotFound) {
	// 	return errs.WithCode(err, errs.ServerErrorCode)
	// }
	//
	// // 配置为空 || 状态为下线 || 值为不跳过
	// // 线上不消费此记录，由快照表手动审核完成后执行
	// if dict == nil || dict.Status == 0 || dict.Value != "1" {
	// 	return nil
	// }
	//
	// // 跳过审核
	// if dict.Value == "1" {
	// 	// _ = p.svcCtx.HubLiveModel.UpdateOperateStatus(ctx, aid, 1, 0, "跳过审核")
	// 	// err = p.svcCtx.MQProducer.SendSync(ctx, define.HubConsumerToPro, p.f)
	// 	// if err != nil {
	// 	// 	return errs.WithMsg(err, errs.ErrPushMQ, fmt.Sprintf("发送MQ消息体失败"))
	// 	// }
	// }
	return nil
}

// ProToTest 直接测试环境数据处理
func (p *ProviderLive) ProToTest(ctx context.Context) error {
	// 错误通知
	if err := render.NewRender(p.svcCtx).Render(ctx, p.f); err != nil {
		return err
	}

	// 推送消息
	userId := p.getUserId(ctx)
	if userId != "" {
		_ = p.PushMessageByCard(ctx, userId)
	} else {
		_ = p.PushMessageToAdmin(ctx)
	}

	return nil
}

// ProToPre 直接预发环境数据处理
func (p *ProviderLive) ProToPre(ctx context.Context) error {
	// 错误通知
	if err := render.NewRender(p.svcCtx).Render(ctx, p.f); err != nil {
		return err
	}

	// 推送消息
	userId := p.getUserId(ctx)
	if userId != "" {
		_ = p.PushMessageByCard(ctx, userId)
	} else {
		_ = p.PushMessageToAdmin(ctx)
	}

	return nil
}

// ProToTestPass 测试环境无误，准备消费线上数据
func (p *ProviderLive) ProToTestPass(ctx context.Context) error {
	return render.NewRender(p.svcCtx).Render(ctx, p.f)
}

// 获取分布式锁
func (p *ProviderLive) getLock(ctx context.Context, key string, value string) error {
	// 尝试获取分布式锁
	ok, err := p.svcCtx.RedisClient.SetnxEx(ctx, key, value, liveCourseLockExpire)
	if err != nil {
		return err
	}
	if !ok {
		logz.Warnf(ctx, "直播课堂，请求频繁，获取锁失败，key=%s,value=%s", key, value)
		return errs.NewCode(errs.ErrTooFast)
	}

	return nil
}

func (p *ProviderLive) getUserId(ctx context.Context) string {
	// 获取用户user_id
	user, err := p.svcCtx.SysUserModel.FindOneByAccount(ctx, p.f.Account)
	if err != nil {
		logc.Errorf(ctx, "查询用户异常，Err:%v", err)
		return ""
		// // 若推送用户为空，则推送消息
		// if errors.Is(err, sql.ErrNoRows) {
		// 	p.svcCtx.SysUserModel.Insert(ctx, &workbench.User{
		// 		Account:      p.f.Account,
		// 		Password:     util.GenerateMD5Str(define.SysNewUserDefaultPassword + p.svcCtx.Config.Salt),
		// 		Username:     p.f.Account,
		// 		Mobile:       "",
		// 		OpenId:       "",
		// 		UserId:       "",
		// 		Email:        "",
		// 		Nickname:     "",
		// 		Gender:       0,
		// 		Avatar:       "",
		// 		ProfessionId: 0,
		// 		JobId:        0,
		// 		DeptId:       0,
		// 		RoleIds:      "",
		// 		Status:       0,
		// 	})
		// } else {
		// logc.Errorf(ctx, "查询用户异常，Err:%v", err)
		// }
	}
	modifyUser := &workbench.User{}
	if user.OpenId == "" {
		// 获取用户OpenId
		openIdMap := feishu.GetOpenId(user.Mobile, "open_id")
		if openId, ok := openIdMap[user.Mobile]; ok {
			user.OpenId, modifyUser.OpenId = openId, openId
		}
	}
	if user.UserId == "" {
		// 获取用户UserId
		userIdMap := feishu.GetOpenId(user.Mobile, "user_id")
		if userId, ok := userIdMap[user.Mobile]; ok {
			user.UserId, modifyUser.UserId = userId, userId
		}
	}
	// 更新表
	if modifyUser.OpenId != "" || modifyUser.UserId != "" {
		_, _ = p.svcCtx.SysUserModel.UpdateFillFieldsById(ctx, user.Id, modifyUser)
	}

	return user.UserId
}

func (p *ProviderLive) PushMessageByCard(ctx context.Context, userId string) error {
	// 消费成功通知
	receive := feishu.UserReceive{
		UserType: "user_id",
		UserId:   userId,
	}
	content := feishu.TemplateVariable{
		Source:    define.HubSourceMap[p.f.Source],
		Title:     p.f.Name,
		Mode:      define.ModeNameMap[p.svcCtx.Config.Mode],
		Operation: define.LiveOperateMap[p.f.EnvType],
		Number:    p.f.LiveNo,
		At: []string{
			receive.UserId,
		},
		Topic: p.svcCtx.Config.Cls.TopicID,
	}
	err := feishu.SendCardMessageByApp(ctx, receive, feishu.CardTemplateSuccess, content)
	if err != nil {
		logc.Errorf(ctx, "飞书推送消息通知有误，Err:%v", err)
	}

	return nil
}

func (p *ProviderLive) PushMessageToAdmin(ctx context.Context) error {
	// 消费成功通知
	receive := feishu.UserReceive{
		UserType: "user_id",
		UserId:   "3f2a1d64", // 王鹏飞
	}
	// 腾讯云访问日志的内容base64编码
	traceId := trace.TraceIDFromContext(ctx)
	if traceId == "" {
		traceId = "暂无"
	}
	encoded := base64.StdEncoding.EncodeToString([]byte(traceId))
	content := feishu.RichText{
		ZhCn: feishu.RichTextLanguageContent{
			Title: fmt.Sprintf("推送人：%s", p.f.Account),
			Content: [][]feishu.RichTextContentElement{
				{{Tag: "text", Text: fmt.Sprintf("数据编号：%s", p.f.LiveNo)}},
				{{Tag: "text", Text: fmt.Sprintf("数据标题：%s", p.f.Name)}},
				{
					{Tag: "text", Text: "追踪记录："},
					{Tag: "a", Text: traceId, Href: "https://console.cloud.tencent.com/cls/search?region=ap-beijing&topic_id=" + p.svcCtx.Config.Cls.TopicID + "&queryBase64=" + encoded + "&time=now-h,now"},
				},
				{
					// {Tag: "text", Text: "请求结果：推送成功，但无指定人ID导致无法推送消息，请及时添加", Style: []string{"bold"}},
					{Tag: "text", Text: "请求结果：推送成功，但无指定人ID导致无法推送消息，请及时添加"},
				},
			},
		},
	}
	err := feishu.SendRichTextMessageByApp(ctx, receive, content)
	if err != nil {
		logc.Errorf(ctx, "飞书推送消息通知有误，Err:%v", err)
	}

	return nil
}
