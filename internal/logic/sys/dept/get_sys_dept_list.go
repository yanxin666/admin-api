package dept

import (
	"context"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetSysDeptListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSysDeptListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSysDeptListLogic {
	return &GetSysDeptListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSysDeptListLogic) GetSysDeptList() (resp *types.SysDeptListResp, err error) {
	sysDeptList, err := l.svcCtx.SysDeptModel.FindAll(l.ctx)
	if err != nil {
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}

	var dept types.Dept
	deptList := make([]types.Dept, 0)
	for _, v := range sysDeptList {
		err := copier.Copy(&dept, &v)
		if err != nil {
			return nil, errs.WithCode(err, errs.ServerErrorCode)
		}
		deptList = append(deptList, dept)
	}

	return &types.SysDeptListResp{
		List: deptList,
	}, nil
}
