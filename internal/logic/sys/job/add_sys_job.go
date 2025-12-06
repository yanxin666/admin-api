package job

import (
	"context"
	"errors"
	"muse-admin/internal/model/workbench"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddSysJobLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddSysJobLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddSysJobLogic {
	return &AddSysJobLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddSysJobLogic) AddSysJob(req *types.AddSysJobReq) error {
	_, err := l.svcCtx.SysJobModel.FindOneByName(l.ctx, req.Name)
	if errors.Is(err, workbench.ErrNotFound) {
		var sysJob = new(workbench.Job)
		err = copier.Copy(sysJob, req)
		if err != nil {
			return errs.WithCode(err, errs.ServerErrorCode)
		}
		_, err = l.svcCtx.SysJobModel.Insert(l.ctx, sysJob)
		if err != nil {
			return errs.WithCode(err, errs.ServerErrorCode)
		}

		return nil
	} else {

		return errs.NewCode(errs.AddJobErrorCode)
	}
}
