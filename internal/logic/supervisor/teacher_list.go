package supervisor

import (
	"context"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TeacherList struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTeacherList(ctx context.Context, svcCtx *svc.ServiceContext) *TeacherList {
	return &TeacherList{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TeacherList) TeacherList(req *types.TeacherListReq) (resp *types.TeacherListResp, err error) {

	condition := make(map[string]interface{})
	if req.Role != 0 {
		condition["role"] = req.Role
	}

	teacherData, err := l.svcCtx.DsTeacherModel.FindAllRole(l.ctx, condition)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	resp = &types.TeacherListResp{List: make([]types.TeacherInfo, 0)}
	for _, item := range teacherData {
		resp.List = append(resp.List, types.TeacherInfo{
			Id:        item.Id,
			Name:      item.Name,
			AccountId: item.AccountId,
			Phone:     item.Phone,
		})
	}

	return
}
