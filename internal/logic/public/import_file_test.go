package public

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"muse-admin/internal/svc"
	"muse-admin/internal/svc/public"
	"muse-admin/internal/types"
	"testing"
	"time"
)

func TestImportFile_ImportFile(t *testing.T) {
	ctx, svcCtx, logger := svc.InitTest()
	type fields struct {
		Logger       logx.Logger
		ctx          context.Context
		svcCtx       *svc.ServiceContext
		alarmService *public.BizAlarm
	}
	type args struct {
		req *types.ImportFileReq
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantResp *types.Result
		wantErr  bool
	}{
		{
			name: "",
			fields: fields{
				Logger:       logger,
				ctx:          ctx,
				svcCtx:       svcCtx,
				alarmService: public.NewAlarmService(svcCtx),
			},
			args: args{
				req: &types.ImportFileReq{},
			},
			wantResp: nil,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &ImportFile{
				Logger:       tt.fields.Logger,
				ctx:          tt.fields.ctx,
				svcCtx:       tt.fields.svcCtx,
				alarmService: tt.fields.alarmService,
			}
			_, _ = l.ImportFile(tt.args.req)

			time.Sleep(10 * time.Second)
			//// 模拟并发
			//var wg sync.WaitGroup
			//for i := 0; i < 10; i++ {
			//	wg.Add(1)
			//	go func(n int) {
			//		defer wg.Done()
			//		_, _ = k.Explain(tt.args.ctx, tt.args.dataStr)
			//	}(i)
			//}
			//wg.Wait() // 等待所有 goroutine 完成
		})
	}
}
