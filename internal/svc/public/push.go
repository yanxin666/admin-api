package public

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/third/feishu"
	"encoding/base64"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/trace"
	"muse-admin/internal/model/workbench"
	"muse-admin/internal/svc"
)

type PushMsgData struct {
	Source    string `json:"source"`    // 数据来源 define.HubSourceMap
	Mode      string `json:"mode"`      // 执行环境 define.ModeNameMap
	Operation string `json:"operation"` // 执行操作 define.LiveOperateMap
	Title     string `json:"title"`     // 数据标题
	Number    string `json:"number"`    // 数据编号
	Account   string `json:"account"`   // 操作者账号
}

type BizPushMsg struct {
	svcCtx *svc.ServiceContext
}

func NewBizPushMsgService(svcCtx *svc.ServiceContext) *BizPushMsg {
	return &BizPushMsg{
		svcCtx: svcCtx,
	}
}

// PushMsgCard 推送卡片消息 将 internal/consumer/hub/live/provider.go line 147~152 中的推送功能独立出来
func (l *BizPushMsg) PushMsgCard(ctx context.Context, data PushMsgData) error {
	// 获取用户user_id
	userId := l.getUserId(ctx, l.svcCtx, data.Account)

	// 向操作者推送卡片信息
	if userId != "" {
		return l.pushMessageByCard(ctx, userId, data)
	}

	// 向管理员推送消息
	return l.pushMessageToAdmin(ctx, data)
}

func (l *BizPushMsg) getUserId(ctx context.Context, svcCtx *svc.ServiceContext, account string) string {
	if account == "" {
		return ""
	}

	// 获取用户user_id
	user, err := svcCtx.SysUserModel.FindOneByAccount(ctx, account)
	if err != nil {
		logc.Errorf(ctx, "查询用户异常，Err:%v", err)
		return ""
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
		_, _ = svcCtx.SysUserModel.UpdateFillFieldsById(ctx, user.Id, modifyUser)
	}

	return user.UserId
}

func (l *BizPushMsg) pushMessageByCard(ctx context.Context, userId string, data PushMsgData) error {
	// 消费成功通知
	receive := feishu.UserReceive{
		UserType: "user_id",
		UserId:   userId,
	}
	content := feishu.TemplateVariable{
		Source:    data.Source,
		Title:     data.Title,
		Mode:      data.Mode,
		Operation: data.Operation,
		Number:    data.Number,
		At: []string{
			receive.UserId,
		},
		Topic: l.svcCtx.Config.Cls.TopicID,
	}
	err := feishu.SendCardMessageByApp(ctx, receive, feishu.CardTemplateSuccess, content)
	if err != nil {
		logc.Errorf(ctx, "飞书推送消息通知有误，Err:%v", err)
	}

	return nil
}

func (l *BizPushMsg) pushMessageToAdmin(ctx context.Context, data PushMsgData) error {
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
			Title: fmt.Sprintf("推送人：%s", data.Account),
			Content: [][]feishu.RichTextContentElement{
				{{Tag: "text", Text: fmt.Sprintf("数据编号：%s", data.Number)}},
				{{Tag: "text", Text: fmt.Sprintf("数据标题：%s", data.Title)}},
				{
					{Tag: "text", Text: "追踪记录："},
					{Tag: "a", Text: traceId, Href: "https://console.cloud.tencent.com/cls/search?region=ap-beijing&topic_id=" + l.svcCtx.Config.Cls.TopicID + "&queryBase64=" + encoded + "&time=now-h,now"},
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
