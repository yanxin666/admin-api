package supertrain

import (
	"context"
	"muse-admin/internal/model/workbench"
	"muse-admin/pkg/errs"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type Detail struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetail(ctx context.Context, svcCtx *svc.ServiceContext) *Detail {
	return &Detail{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Detail) Detail(req *types.HubSuperTrainDetailReq) (resp *types.HubSuperTrainInfo, err error) {
	data, err := l.svcCtx.HubSuperTrainModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}

	var userInfo *workbench.User
	operateName := "暂无"
	if data.OperateId != 0 {
		userInfo, _ = l.svcCtx.SysUserModel.FindOne(l.ctx, data.OperateId)
		if userInfo != nil {
			operateName = userInfo.Username
		}
	}

	return &types.HubSuperTrainInfo{
		Id:            data.Id,
		No:            data.No,
		Name:          data.Name,
		Status:        data.Status,
		Version:       data.Version,
		OperateStatus: data.OperateStatus,
		OperateName:   operateName,
		Data:          data.Data,
		AppVersion:    data.AppVersion,
		CreatedAt:     data.CreatedAt.Unix(),
		UpdatedAt:     data.UpdatedAt.Unix(),
	}, nil
}
