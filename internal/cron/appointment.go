package cron

import (
	// appointmentModel "muse-admin/internal/model/appointment"
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"muse-admin/internal/svc"
)

type AppointmentCron struct {
	// 基础数据
	ctx    context.Context
	svcCtx *svc.ServiceContext
	// appointUserModel appointmentModel.UserAppointmentModel
	// 数据存储
}

func NewAppointmentCron(ctx context.Context, svcCtx *svc.ServiceContext) *AppointmentCron {
	return &AppointmentCron{
		ctx:    ctx,
		svcCtx: svcCtx,
		// appointUserModel: appointmentModel.NewUserAppointmentModel(svcCtx.MysqlConn),
	}
}

// Access 是否准入
func (c *AppointmentCron) Access() bool {
	return false
}

// SetTitle 脚本名称
func (c *AppointmentCron) SetTitle() string {
	return "定时任务-更新用户预约课程脚本"
}

func (c *AppointmentCron) SetCron() string {
	// */30 表示分钟部分，表示每 30 分钟执行一次任务
	return "0 */30 * * * *"
}

func (c *AppointmentCron) Run() {
	// 这里是定时任务执行的具体逻辑

	fmt.Println(c.SetTitle()+"已执行，当前时间：", util.GetStandardNowDatetime())

	// 获取当前时间
	time := util.GetStandardNowDatetime()
	fmt.Println(time)

	// 每半小时刷新一下用户预约课程的状态
	// err := c.appointUserModel.CronUpdateUserStatus(c.ctx, time, &appointmentModel.UserAppointment{
	// 	Status: define.AppointmentUserStatus.Expired,
	// })

	// // 脚本异常报警通知
	// if err != nil {
	// 	content := "注意：执行" + c.SetTitle() + "失败！请及时检查,Err:%v"
	// 	// 记录日志
	// 	logc.Errorf(c.ctx, content, err)
	// 	// 发送飞书报警
	// 	_ = public.NewAlarmService(c.ctx, c.svcCtx).CronAlarm(define.Alarm.Wpf, content, err)
	// 	return
	// }

	// 定时任务执行记录
	logc.Infof(c.ctx, "%s成功！当前时间：%s", c.SetTitle(), util.GetStandardNowDatetime())

	return
}
