package role

import (
	"context"
	"encoding/json"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetSysRoleListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSysRoleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSysRoleListLogic {
	return &GetSysRoleListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSysRoleListLogic) GetSysRoleList() (resp *types.SysRoleListResp, err error) {
	sysRoleList, err := l.svcCtx.SysRoleModel.FindAll(l.ctx)
	if err != nil {
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}

	var role types.Role
	roleList := make([]types.Role, 0)
	for _, v := range sysRoleList {
		err := copier.Copy(&role, &v)
		if err != nil {
			return nil, errs.WithCode(err, errs.ServerErrorCode)
		}
		var permMenuIds []int64
		err = json.Unmarshal([]byte(v.PermMenuIds), &permMenuIds)
		role.PermMenuIds = permMenuIds
		roleList = append(roleList, role)
	}

	return &types.SysRoleListResp{
		List: roleList,
	}, nil
}
