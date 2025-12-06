package session

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"testing"
)

func TestGetListByUserId_GetListByUserId(t *testing.T) {
	ctx, svcCtx, logger := svc.InitTest()
	// 入参测试
	in := &types.GetListByUserIdReq{
		UserId: 1195,
	}
	type fields struct {
		Logger logx.Logger
		ctx    context.Context
		svcCtx *svc.ServiceContext
	}
	type args struct {
		req *types.GetListByUserIdReq
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantResp *types.GetListByUserIdResp
		wantErr  bool
	}{
		{
			name: "根据userid查询会话列表",
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
			l := &GetListByUserId{
				Logger: tt.fields.Logger,
				ctx:    tt.fields.ctx,
				svcCtx: tt.fields.svcCtx,
			}
			gotResp, err := l.GetListByUserId(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetListByUserId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("根据userid查询会话列表-单测结果: %v", gotResp)
		})
	}
}
