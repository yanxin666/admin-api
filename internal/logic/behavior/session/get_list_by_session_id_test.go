package session

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"testing"
)

func TestGetListBySessionId_GetListBySessionId(t *testing.T) {
	ctx, svcCtx, logger := svc.InitTest()
	// 入参测试
	in := &types.GetListBySessionIdReq{
		SessionId: 2,
	}
	type fields struct {
		Logger logx.Logger
		ctx    context.Context
		svcCtx *svc.ServiceContext
	}
	type args struct {
		req *types.GetListBySessionIdReq
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantResp *types.GetListBySessionIdResp
		wantErr  bool
	}{
		{
			name: "根据sessionid查询会话记录列表",
			fields: fields{
				Logger: logger,
				ctx:    ctx,
				svcCtx: svcCtx,
			},
			args: args{
				req: in,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &GetListBySessionId{
				Logger: tt.fields.Logger,
				ctx:    tt.fields.ctx,
				svcCtx: tt.fields.svcCtx,
			}
			gotResp, err := l.GetListBySessionId(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetListBySessionId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("根据sessionid查询会话记录列表-单测结果: %v", gotResp)
		})
	}
}
