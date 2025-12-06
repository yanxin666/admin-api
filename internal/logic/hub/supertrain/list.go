package supertrain

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"muse-admin/internal/model/tools"
	"muse-admin/pkg/errs"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type List struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewList(ctx context.Context, svcCtx *svc.ServiceContext) *List {
	return &List{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *List) List(req *types.HubSuperTrainListReq) (resp *types.HubSuperTrainListResp, err error) {
	var (
		keywordNo, keywordName string
	)

	resp = &types.HubSuperTrainListResp{
		List: make([]types.HubSuperTrainInfo, 0),
	}

	// 构建查询条件
	condition := tools.FilterConditions(req)
	delete(condition, "no")
	delete(condition, "name")
	if req.No != "" {
		keywordNo = req.No
	}
	if req.Name != "" {
		keywordName = req.Name
	}

	list, total, err := l.svcCtx.HubSuperTrainModel.FindPageByCondition(l.ctx, req.Page, req.Limit, keywordNo, keywordName, condition)
	if err != nil {
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}
	if len(list) == 0 {
		return resp, nil
	}

	operateIds, _ := util.ArrayColumn(list, "OperateId")
	operateIds = util.ArrayUniqueValue(operateIds)
	userInfo, err := l.svcCtx.SysUserModel.BatchByUserIds(l.ctx, operateIds)
	if err != nil {
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}

	inspectMap, err := l.svcCtx.HubSuperTrainModel.FindMaxVersion(l.ctx)
	if err != nil {
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}

	for _, v := range list {
		// 检查当前数据是否为多版本，若是，但值不等于最新的，代表多版本的旧数据，需要打标签为被覆盖
		if targetValue, exists := inspectMap[v.No]; exists && v.Version != targetValue {
			v.OperateStatus = 99
		}

		resp.List = append(resp.List, types.HubSuperTrainInfo{
			Id:            v.Id,
			No:            v.No,
			Name:          v.Name,
			Status:        v.Status,
			Version:       v.Version,
			OperateStatus: v.OperateStatus,
			OperateName:   userInfo[v.OperateId].Username,
			Data:          "",
			AppVersion:    v.AppVersion,
			CreatedAt:     v.CreatedAt.Unix(),
			UpdatedAt:     v.UpdatedAt.Unix(),
		})
	}

	// 分页信息
	resp.Pagination = types.Pagination{
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	}

	return
}
