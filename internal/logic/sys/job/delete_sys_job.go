package job

import (
	"context"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSysJobLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteSysJobLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSysJobLogic {
	return &DeleteSysJobLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSysJobLogic) DeleteSysJob(req *types.DeleteSysJobReq) error {
	count, _ := l.svcCtx.SysUserModel.FindCountByCondition(l.ctx, "job_id", req.Id)
	if count != 0 {
		return errs.NewCode(errs.DeleteJobErrorCode)
	}

	err := l.svcCtx.SysJobModel.Delete(l.ctx, req.Id)
	if err != nil {
		return errs.WithCode(err, errs.ServerErrorCode)
	}

	return nil
}
