package job

import (
	"context"
	"muse-admin/internal/define"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSysJobLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSysJobLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSysJobLogic {
	return &UpdateSysJobLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSysJobLogic) UpdateSysJob(req *types.UpdateSysJobReq) error {
	sysJob, err := l.svcCtx.SysJobModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return errs.NewCode(errs.JobIdErrorCode)
	}

	if req.Status == define.SysDisable {
		count, _ := l.svcCtx.SysUserModel.FindCountByJobId(l.ctx, req.Id)
		if count > 0 {
			return errs.NewCode(errs.JobIsUsingErrorCode)
		}
	}

	err = copier.Copy(sysJob, req)
	if err != nil {
		return errs.WithCode(err, errs.ServerErrorCode)
	}

	err = l.svcCtx.SysJobModel.Update(l.ctx, sysJob)
	if err != nil {
		return errs.WithCode(err, errs.ServerErrorCode)
	}

	return nil
}
