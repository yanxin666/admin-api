package public

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/third/feishu"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/trace"
	"muse-admin/internal/define"
	"muse-admin/internal/svc"
	"muse-admin/internal/tools"
	"muse-admin/pkg/errs"
)

type BizAlarm struct {
	svcCtx *svc.ServiceContext
}

func NewAlarmService(svcCtx *svc.ServiceContext) *BizAlarm {
	return &BizAlarm{
		svcCtx: svcCtx,
	}
}

// ImportAlarm 数据导入报警机器人
func (l *BizAlarm) ImportAlarm(ctx context.Context, operateId int64, title string, err error) error {
	// 获取操作人信息
	sysUser, _ := l.svcCtx.SysUserModel.FindOne(ctx, operateId)

	// 获取发送对象
	var (
		at       string
		username string
	)
	if sysUser != nil {
		at = sysUser.OpenId
		username = sysUser.Username
	} else {
		at = "all"
		username = "所有人"
	}

	// 若当前操作人暂未有openId，根据手机号去获取openId
	if at == "" {
		// 判断操作人是否填充了手机号
		if sysUser.Mobile != "" {
			userMap := feishu.GetOpenId(sysUser.Mobile, "openId")
			if openId, ok := userMap[sysUser.Mobile]; ok {
				sysUser.OpenId = openId
				_ = l.svcCtx.SysUserModel.Update(ctx, sysUser)
				at = openId
			}
		} else {
			at = "all" // 手机号也没有，需要报警所有人
		}
	}

	// 非预发布和线上环境，不发送消息
	// if l.svcCtx.Config.Mode != "pre" && l.svcCtx.Config.Mode != "pro" {
	//	return nil
	// }
	modeName := define.ModeNameMap[l.svcCtx.Config.Mode]
	content := fmt.Sprintf(`%s
执行环境: %s
执行操作人: %s
错误信息: %s
追踪记录: %s`, title, modeName, username, err, trace.TraceIDFromContext(ctx))
	logc.Info(ctx, content)
	botClient := feishu.NewBotMessage(define.BotDomain + define.ImportAlarm)
	err = botClient.SendTextMsg(ctx, content, []feishu.At{{ID: at}})
	if err != nil {
		return errs.WithErr(err)
	}

	return nil
}

// CronAlarm 定时任务报警机器人
func (l *BizAlarm) CronAlarm(ctx context.Context, id, content string, err error) error {
	botClient := feishu.NewBotMessage(define.BotDomain + define.BotAlarm)
	err = botClient.SendTextMsg(ctx, fmt.Sprintf(content, err), []feishu.At{{ID: id}})
	if err != nil {
		return errs.WithErr(err)
	}

	return nil
}

// BusinessAlarm AI业务报警机器人
func (l *BizAlarm) BusinessAlarm(ctx context.Context, id, title string, err error) error {
	// 非预发布和线上环境，不发送消息 // todo 这里上线记得打开，谁看到记得提醒一下
	// if l.svcCtx.Config.Mode != "pre" && l.svcCtx.Config.Mode != "pro" {
	// 	return nil
	// }

	content := fmt.Sprintf(`%s
Err: %s
Uid:%d
TraceId: %s`, title, err, tools.GetUserIdByCtx(ctx), trace.TraceIDFromContext(ctx))
	botClient := feishu.NewBotMessage(define.BotDomain + define.BotAlarm)
	err = botClient.SendTextMsg(ctx, content, []feishu.At{{ID: id}})
	if err != nil {
		return errs.WithErr(err)
	}

	return nil
}

// 后续别的业务报警，请在此处添加，不建议把bot当成参数吐出，而是每个报警线都需要有自己的方法
// ...
