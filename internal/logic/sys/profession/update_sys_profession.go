package profession

import (
	"context"
	"muse-admin/internal/define"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSysProfessionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSysProfessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSysProfessionLogic {
	return &UpdateSysProfessionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSysProfessionLogic) UpdateSysProfession(req *types.UpdateSysProfessionReq) error {
	sysProfession, err := l.svcCtx.SysProfessionModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return errs.NewCode(errs.ProfessionIdErrorCode)
	}

	if req.Status == define.SysDisable {
		count, _ := l.svcCtx.SysUserModel.FindCountByProfessionId(l.ctx, req.Id)
		if count > 0 {
			return errs.NewCode(errs.JobIsUsingErrorCode)
		}
	}

	err = copier.Copy(sysProfession, req)
	if err != nil {
		return errs.WithCode(err, errs.ServerErrorCode)
	}

	err = l.svcCtx.SysProfessionModel.Update(l.ctx, sysProfession)
	if err != nil {
		return errs.WithCode(err, errs.ServerErrorCode)
	}

	return nil
}
