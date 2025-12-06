package define

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/third/feishu"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"encoding/base64"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/trace"
)

var LiveOperateMap = map[int64]string{
	1: "预览",
	2: "发布",
}

var ModeNameMap = map[string]string{
	"dev":  "开发环境「可忽略」",
	"test": "测试环境",
	"pro":  "线上环境",
}

const (
	Lesson     int64 = iota + 1 // 统编教材精讲课堂
	WritePPT                    // 写作PPT
	AILive                      // AI直播课堂
	SuperTrain                  // 超能训练
)

var HubSourceMap = map[int64]string{
	Lesson:     "统编教材精讲课堂",
	WritePPT:   "写作PPT",
	AILive:     "AI直播课堂",
	SuperTrain: "超能训练",
}

// 机器人消息通知
const (
	// BotDomain 飞书机器人请求地址
	BotDomain = "https://open.feishu.cn/open-apis/bot/v2/hook/"

	// BotAlarm 定时任务报警机器人
	BotAlarm = "fe2254c9-411b-4cd1-835d-6a50530cc9e6"

	// MQAlarm MQ消费报警机器人
	MQAlarm = "946533e7-6898-47e4-8b5e-6fe2a55832f0"

	// ImportAlarm 数据导入报警机器人
	ImportAlarm = "c5adef64-8d65-4213-a149-fd778d14fce2"

	// Self 自己测试使用
	Self = "d4c5adc9-c255-4226-99a3-fca476f5d129"
)

var AlarmNotify = map[string]string{
	"":    "all",
	"王鹏飞": "ou_0dfb63804fbab43ddbca3d0987b77b41",
	"王林":  "ou_fbd456786feafab28a03cff292701319",
	"严鑫":  "ou_7c32654398dc3f231f6900b2cad7a77a",
	"张兴":  "ou_8e7a9f401ceeff594818fc9db4de543e",
	"魏帅":  "ou_05ac63f14aed747f44ab8bdb93bb6bbb",
	"王玉妹": "ou_ce921c0ea29c26a932cb21f17d2e06aa",
	"张金涛": "ou_3adf9a7c496105e9aebb1fae52233b40",
	"王尚峰": "ou_2a042e1ac8f92d76d0a5f7b3024a28ca",
}

// Alarm 监控报警通知人员
var Alarm = struct {
	All string
	Yx  string
	Wpf string
	Zx  string
	Wl  string
}{
	All: "all",
	Yx:  "ou_7c32654398dc3f231f6900b2cad7a77a",
	Wpf: "ou_0dfb63804fbab43ddbca3d0987b77b41",
	Zx:  "ou_8e7a9f401ceeff594818fc9db4de543e",
	Wl:  "ou_fbd456786feafab28a03cff292701319",
}

func AlarmMqCustomer(ctx context.Context, mode, id, clsTopicId, title, topic, body string, err error) error {
	// 非预发布和线上环境，不发送消息
	// if mode != "pre" && mode != "pro" {
	//	return nil
	// }
	if topic == "" {
		topic = "暂无"
	}

	// 记录完整数据
	logc.Info(ctx, body)

	// 确保字符串长度不超过250
	body = util.SafeSubstring(body, 250)

	// 腾讯云访问日志的内容base64编码
	traceId := trace.TraceIDFromContext(ctx)
	encoded := base64.StdEncoding.EncodeToString([]byte(traceId))
	modeName := ModeNameMap[mode]

	content := feishu.RichTextContent{
		Title: title,
		Contents: [][]feishu.RichContentItem{
			{{Tag: "text", Text: "操作人员："}},
			{{Tag: "text", Text: fmt.Sprintf("主题名称：%s", topic)}},
			{{Tag: "text", Text: fmt.Sprintf("执行环境：%s", modeName)}},
			{{Tag: "text", Text: fmt.Sprintf("错误信息：%s", err)}},
			{{Tag: "text", Text: fmt.Sprintf("推送数据：%s", body)}},
			{{Tag: "text", Text: "追踪记录："}, {Tag: "a", Text: traceId, Href: "https://console.cloud.tencent.com/cls/search?region=ap-beijing&topic_id=" + clsTopicId + "&queryBase64=" + encoded + "&time=now-h,now"}},
		},
	}
	botClient := feishu.NewBotMessage(BotDomain + MQAlarm)
	err = botClient.SendRichTextMessage(ctx, content, []feishu.At{{ID: id}})
	if err != nil {
		return err
	}

	return nil
}
