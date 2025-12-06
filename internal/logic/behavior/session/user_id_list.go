package session

import (
	"context"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserIdList struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserIdList(ctx context.Context, svcCtx *svc.ServiceContext) *UserIdList {
	return &UserIdList{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserIdList) UserIdList() (resp *types.UserIdListResp, err error) {
	// 查询数据
	list, err := l.svcCtx.SessionUserModel.FindUidList(l.ctx)
	if err != nil {
		return nil, err
	}
	// 暂无数据
	if len(list) == 0 {
		return nil, nil
	}

	// 构造返回结果
	return &types.UserIdListResp{
		List: list,
	}, nil
}
