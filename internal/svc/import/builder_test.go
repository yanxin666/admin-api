package _import

import (
	"context"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"muse-admin/internal/config"
	taskModel "muse-admin/internal/model/task"
	"muse-admin/internal/svc"
	"testing"
)

func InitTest() (context.Context, *svc.ServiceContext, logx.Logger) {
	var c config.Config
	var configFile = "../../../etc/application.yaml"

	conf.MustLoad(configFile, &c)
	ctx := context.Background()
	svcCtx := svc.NewServiceContext(c)
	log := logx.WithContext(ctx)
	return ctx, svcCtx, log
}

func TestNewBuilder(t *testing.T) {
	ctx, svcCtx, _ := InitTest()
	type args struct {
		svcCtx      *svc.ServiceContext
		taskInfo    *taskModel.SyncTask
		filename    string
		operateName string
	}
	tests := []struct {
		name string
		args args
		want *BaseBuilder
	}{
		{
			name: "",
			args: args{
				svcCtx: svcCtx,
				taskInfo: &taskModel.SyncTask{
					Id:        192,
					Type:      1,
					OperateId: 3,
					FileSheet: 0,
				},
				filename:    "/Users/wangpengfei/5734_1755842014617.xlsx",
				operateName: "王鹏飞",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewBuilder(tt.args.svcCtx, tt.args.taskInfo, tt.args.filename).BuildStage(ctx)
		})
	}
}
