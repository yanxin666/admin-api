package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/stat"
	"github.com/zeromicro/go-zero/rest"
	handler2 "github.com/zeromicro/go-zero/rest/handler"
	"muse-admin/internal/config"
	"muse-admin/internal/handler"
	"muse-admin/internal/middleware"
	"muse-admin/internal/svc"
)

var configFile = flag.String("f", "etc/application.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf, rest.WithNotAllowedHandler(middleware.NewCorsMiddleware().Handler()))
	defer server.Stop()

	svcCtx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, svcCtx)
	stat.DisableLog()

	// 全局中间件
	server.Use(middleware.NewRecoveryMiddleware(c).Handle)     // CLS日志服务 & panic处理
	server.Use(middleware.NewRequestLogMiddleware(c).Handle)   // 入参明细记录
	server.Use(middleware.NewCorsMiddleware().Handle)          // 跨域
	server.Use(rest.ToMiddleware(handler2.DetailedLogHandler)) // API响应明细日志

	// 挂起Cron定时任务
	//cron.InitClient(svcCtx)

	// MQ消费注册
	//consumer.Register(c, svcCtx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
