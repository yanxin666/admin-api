package profession

import (
	"context"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSysProfessionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteSysProfessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSysProfessionLogic {
	return &DeleteSysProfessionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSysProfessionLogic) DeleteSysProfession(req *types.DeleteSysProfessionReq) error {
	count, _ := l.svcCtx.SysUserModel.FindCountByCondition(l.ctx, "profession_id", req.Id)
	if count != 0 {
		return errs.NewCode(errs.DeleteProfessionErrorCode)
	}

	err := l.svcCtx.SysProfessionModel.Delete(l.ctx, req.Id)
	if err != nil {
		return errs.WithCode(err, errs.ServerErrorCode)
	}

	return nil
}
