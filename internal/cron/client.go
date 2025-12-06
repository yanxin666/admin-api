package cron

import (
	"context"
	pg "e.coding.net/zmexing/nenglitanzhen/biz-lib/cron"
	"muse-admin/internal/svc"
)

func InitClient(svcCtx *svc.ServiceContext) {
	ctx := context.Background()

	c := pg.NewCrontabClient()

	// 注册预约课程
	// if r, err := NewAppointmentCron(ctx); err != nil {
	// 	c.RegisterCron(r)
	// }

	// 注册预约课程
	if r := NewAppointmentCron(ctx, svcCtx); r != nil {
		// c.RegisterCron(r)
	}

	// 注册自动排课
	if r := NewAutoIncrLessonCorn(ctx, svcCtx); r != nil {
		// c.RegisterCron(r)
	}
	// 触发定时任务
	c.Exec()
}
