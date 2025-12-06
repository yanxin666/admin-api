package order

import (
	"context"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExportOrder struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExportOrder(ctx context.Context, svcCtx *svc.ServiceContext) *ExportOrder {
	return &ExportOrder{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExportOrder) ExportOrder(req *types.GetOrderListReq) (resp *types.GetOrderListResp, err error) {
	if err = l.checkParams(req); err != nil {
		return nil, err
	}

	// todo: add your logic here and delete this line
	// ...

	return
}

// 校验参数
func (l *ExportOrder) checkParams(req *types.GetOrderListReq) error {
	// todo: add your logic check params and delete this line
	// 若校验的参数过少，可删除此方法，直接在上方 ExportOrder 中编写

	return nil
}
