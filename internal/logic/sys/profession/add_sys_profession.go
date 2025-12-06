package profession

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

type AddSysProfessionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddSysProfessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddSysProfessionLogic {
	return &AddSysProfessionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddSysProfessionLogic) AddSysProfession(req *types.AddSysProfessionReq) error {
	_, err := l.svcCtx.SysProfessionModel.FindOneByName(l.ctx, req.Name)
	if errors.Is(err, workbench.ErrNotFound) {
		var sysProfession = new(workbench.Profession)
		err = copier.Copy(sysProfession, req)
		if err != nil {
			return errs.WithCode(err, errs.ServerErrorCode)
		}
		_, err = l.svcCtx.SysProfessionModel.Insert(l.ctx, sysProfession)
		if err != nil {
			return errs.WithCode(err, errs.ServerErrorCode)
		}

		return nil
	} else {

		return errs.NewCode(errs.AddProfessionErrorCode)
	}
}
