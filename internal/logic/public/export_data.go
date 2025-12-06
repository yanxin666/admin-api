package public

import (
	"context"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExportData struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExportData(ctx context.Context, svcCtx *svc.ServiceContext) *ExportData {
	return &ExportData{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExportData) ExportData(req *types.ExportDataReq) (resp *types.Result, err error) {
	if err = l.checkParams(req); err != nil {
		return nil, err
	}

	// todo: add your logic here and delete this line
	// ...

	return
}

// 校验参数
func (l *ExportData) checkParams(req *types.ExportDataReq) error {
	// todo: add your logic check params and delete this line
	// 若校验的参数过少，可删除此方法，直接在上方 ExportData 中编写

	return nil
}
