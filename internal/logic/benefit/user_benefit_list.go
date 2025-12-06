package benefit

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logc"
	"muse-admin/internal/model/benefit"
	"muse-admin/pkg/errs"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserBenefitList struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserBenefitList(ctx context.Context, svcCtx *svc.ServiceContext) *UserBenefitList {
	return &UserBenefitList{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserBenefitList) UserBenefitList(req *types.UserBenefitListReq) (resp *types.UserBenefitListResp, err error) {
	resp = &types.UserBenefitListResp{
		List: make([]types.UserBenefitInfo, 0),
	}

	var (
		keyword string
		userIds []any
		list    []benefit.UserBenefitsDetail
		total   int64
	)

	// 构建查询条件
	conditionByUser := make(map[string]interface{})
	if req.UserId != 0 {
		userIds = append(userIds, req.UserId)
	}
	if len(req.Phone) == 11 {
		// 加密手机号
		conditionByUser["mask_phone"], err = util.AesEncrypt(l.ctx, req.Phone, l.svcCtx.Config.EncryptKey)
		if err != nil {
			return nil, err
		}
	} else {
		keyword = req.Phone
	}

	if len(conditionByUser) != 0 || keyword != "" {
		// 获取主用户列表
		baseList, err := l.svcCtx.BaseUserModel.FindAllByCondition(l.ctx, keyword, conditionByUser)
		if err != nil {
			return nil, errs.WithCode(err, errs.ServerErrorCode)
		}
		// 无用户
		if baseList == nil {
			return resp, nil
		}
		baseUserIds, _ := util.ArrayColumn(baseList, "Id")
		baseUserIds = util.ArrayUniqueValue(baseUserIds)
		// 获取符合条件的所有子用户
		userMap, _ := l.svcCtx.UserModel.BatchMapByBaseUserIds(l.ctx, baseUserIds)
		for _, info := range userMap {
			for _, v := range info {
				userIds = append(userIds, v.Id)
			}
		}
	}

	// 构建查询条件
	condition := make(map[string]interface{})
	if req.Status != 0 {
		condition["status"] = req.Status
	}
	if req.BenefitsResourceId != 0 {
		condition["benefits_resource_id"] = req.BenefitsResourceId
	}

	if len(userIds) == 0 {
		list, total, err = l.svcCtx.UserBenefitModel.FindPageByCondition(l.ctx, req.Page, req.Limit, condition, req.FromTime, req.EndTime)
	} else {
		list, total, err = l.svcCtx.UserBenefitModel.FindPageByConditionAndUserIds(l.ctx, req.Page, req.Limit, userIds, condition, req.FromTime, req.EndTime)
	}
	if err != nil {
		logc.Errorf(l.ctx, "渠道订单关联渠道子订单查询失败，Err:%s", err)
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}

	if len(list) == 0 {
		return resp, nil
	}

	userIdArr, _ := util.ArrayColumn(list, "UserId")
	userIdArr = util.ArrayUniqueValue(userIdArr)
	// 批量用户信息
	batchUserMap, err := l.svcCtx.UserModel.BatchByUserIds(l.ctx, userIdArr)
	if err != nil {
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}

	resourceIdArr, _ := util.ArrayColumn(list, "BenefitsResourceId")
	resourceIdArr = util.ArrayUniqueValue(resourceIdArr)
	// 批量权益详情信息
	resourceMap, err := l.svcCtx.BenefitResourceModel.BatchByIds(l.ctx, resourceIdArr)
	if err != nil {
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}

	userResourceMap := make(map[int64][]benefit.GroupResourceByUserIdResult)

	// 获取所有用户权益记录
	if len(userIdArr) > 0 {
		for _, userId := range userIdArr {
			uid := cast.ToInt64(userId)
			groupResource, err := l.svcCtx.UserBenefitRecordModel.GroupResourceByUserId(l.ctx, uid)
			if err != nil {
				return nil, errs.WithCode(err, errs.ServerErrorCode)
			}

			userResourceMap[uid] = groupResource
		}
	}

	for _, v := range list {
		item := types.UserBenefitInfo{
			Id:                 v.Id,
			BaseUserId:         batchUserMap[v.UserId].BaseUserId,
			UserNo:             batchUserMap[v.UserId].UserNo,
			Phone:              batchUserMap[v.UserId].Phone,
			UserId:             v.UserId,
			RealName:           batchUserMap[v.UserId].RealName,
			BenefitsResourceId: v.BenefitsResourceId,
			Status:             v.Status,
			BenefitName:        resourceMap[v.BenefitsResourceId].Name,
			StartTime:          v.FromTime.Unix(),
			EndTime:            v.EndTime.Unix(),
			CreatedAt:          v.CreatedAt.Unix(),
			UpdatedAt:          v.UpdatedAt.Unix(),
		}
		// 用户权益发放与回收统计
		if userResourceMap != nil {
			for _, result := range userResourceMap[v.UserId] {
				if result.BenefitsResourceId == v.BenefitsResourceId {
					item.GrantCount = result.GrantCount
					item.ReclaimCount = result.ReclaimCount
				}
			}
		}
		resp.List = append(resp.List, item)
	}

	// 分页信息
	resp.Pagination = types.Pagination{
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	}

	return
}
