package cron

import (
	"context"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"muse-admin/internal/config"
	"muse-admin/internal/svc"
	"testing"
)

func TestAutoIncrLesson(t *testing.T) {
	ctx, svcCtx, _ := testInit()
	corn := NewAutoIncrLessonCorn(ctx, svcCtx)
	corn.Run()
}

// InitTest 单元测试初始化依赖注册
func testInit() (context.Context, *svc.ServiceContext, logx.Logger) {
	var c config.Config
	var configFile = "../../etc/application.yaml"

	conf.MustLoad(configFile, &c)
	ctx := context.Background()
	svcCtx := svc.NewServiceContext(c)
	log := logx.WithContext(ctx)
	return ctx, svcCtx, log
}
