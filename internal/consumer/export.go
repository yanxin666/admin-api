package consumer

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/mq/types"
	rmq "github.com/apache/rocketmq-clients/golang/v5"
	"muse-admin/internal/svc"
)

type Test struct {
	svcCtx *svc.ServiceContext
	a      int
}

func NewTest(svcCtx *svc.ServiceContext) types.Consumer {
	return &Test{
		svcCtx: svcCtx,
	}
}

func (s *Test) Explain(ctx context.Context, data interface{}) error {
	// user, err := s.svcCtx.SysUserModel.FindOne(ctx, 1)
	// fmt.Println(user, err)
	return nil
}

// Execute 区分来源以及去向，角色定位于调度器
func (s *Test) Execute(ctx context.Context) (err error) {
	return nil
}

// ExecuteBatch 区分来源以及去向，角色定位于调度器
func (s *Test) ExecuteBatch(ctx context.Context) ([]*rmq.MessageView, error) {
	return nil, nil
}
