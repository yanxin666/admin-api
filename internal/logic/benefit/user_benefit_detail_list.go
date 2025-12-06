package benefit

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"muse-admin/internal/define"
	"muse-admin/internal/model/benefit"
	"muse-admin/internal/model/tools"
	"muse-admin/pkg/errs"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserBenefitDetailList struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserBenefitDetailList(ctx context.Context, svcCtx *svc.ServiceContext) *UserBenefitDetailList {
	return &UserBenefitDetailList{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserBenefitDetailList) UserBenefitDetailList(req *types.UserBenefitDetailListReq) (resp *types.UserBenefitDetailListResp, err error) {
	resp = &types.UserBenefitDetailListResp{
		List:       make([]types.UserBenefitDetailInfo, 0),
		Pagination: types.Pagination{},
	}

	if req.UserId <= 0 {
		return nil, errs.NewCode(errs.ServerErrorCode)
	}

	// 构建查询条件
	condition := tools.FilterConditions(req)

	// 获取任务详情Log列表
	list, total, err := l.svcCtx.UserBenefitRecordModel.FindPageByCondition(l.ctx, req.Page, req.Limit, condition)
	if err != nil {
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}

	if len(list) == 0 {
		return resp, nil
	}

	resourceIdArr, _ := util.ArrayColumn(list, "BenefitsResourceId")
	resourceIdArr = util.ArrayUniqueValue(resourceIdArr)
	// 批量权益详情信息
	resourceMap, err := l.svcCtx.BenefitResourceModel.BatchByIds(l.ctx, resourceIdArr)
	if err != nil {
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}

	// 按order_id分组，判断是否同时存在type=1和type=2
	orderBool := make(map[string]bool)

	// 按order_id分组
	orderGroups := make(map[string][]benefit.BenefitsUserRecord)
	for _, r := range list {
		orderGroups[r.OrderNo] = append(orderGroups[r.OrderNo], r)
	}

	// 检查每组是否同时有type=1和type=2
	for orderNo, group := range orderGroups {
		hasType1 := false
		hasType2 := false

		for _, r := range group {
			if r.Type == 1 {
				hasType1 = true
			} else if r.Type == 2 {
				hasType2 = true
			}
		}

		// 如果同时存在type=1和type=2，则设置标识
		if hasType1 && hasType2 {
			orderBool[orderNo] = true
		}
	}

	for _, v := range list {
		log := types.UserBenefitDetailInfo{
			Id:            v.Id,
			UserId:        v.UserId,
			Type:          v.Type,
			OrderNo:       v.OrderNo,
			BenefitName:   resourceMap[v.BenefitsResourceId].Name,
			BenefitType:   define.BenefitsTypeMapInt[resourceMap[v.BenefitsResourceId].Type],
			UserBenefitId: v.UserBenefitsId,
			HasBoth:       orderBool[v.OrderNo],
			Source:        define.BenefitsSourceMapInt[v.Source],
			CreatedAt:     v.CreatedAt.Unix(),
		}
		resp.List = append(resp.List, log)
	}

	// 分页信息
	resp.Pagination = types.Pagination{
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	}

	return resp, nil
}
