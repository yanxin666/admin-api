package session

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"testing"
)

func TestUserIdList_UserIdList(t *testing.T) {
	ctx, svcCtx, logger := svc.InitTest()
	type fields struct {
		Logger logx.Logger
		ctx    context.Context
		svcCtx *svc.ServiceContext
	}
	tests := []struct {
		name     string
		fields   fields
		wantResp *types.UserIdListResp
		wantErr  bool
	}{
		{
			name: "获取userId列表-单侧",
			fields: fields{
				Logger: logger,
				ctx:    ctx,
				svcCtx: svcCtx,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &UserIdList{
				Logger: tt.fields.Logger,
				ctx:    tt.fields.ctx,
				svcCtx: tt.fields.svcCtx,
			}
			gotResp, err := l.UserIdList()
			if (err != nil) != tt.wantErr {
				t.Errorf("UserIdList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("获取userId列表单侧结果: %v", gotResp)
		})
	}
}
