package benefit

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"github.com/zeromicro/go-zero/core/logc"
	"muse-admin/internal/define"
	"muse-admin/pkg/errs"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ImportBenefitList struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewImportBenefitList(ctx context.Context, svcCtx *svc.ServiceContext) *ImportBenefitList {
	return &ImportBenefitList{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ImportBenefitList) ImportBenefitList(req *types.ImportBenefitListReq) (resp *types.ImportBenefitListResp, err error) {
	resp = &types.ImportBenefitListResp{
		List:       make([]types.OrderBenefitInfo, 0),
		Pagination: types.Pagination{},
	}

	var keywordOrderNo string
	// 构建查询条件
	condition := make(map[string]interface{})
	if req.Phone != "" {
		// 加密手机号
		phoneEncrypt, err := util.AesEncrypt(l.ctx, req.Phone, l.svcCtx.Config.EncryptKey)
		if err != nil {
			return nil, err
		}
		condition["mask_phone"] = phoneEncrypt
	}
	if req.OrderNo != "" {
		keywordOrderNo = req.OrderNo
	}
	if req.IsGrant != 0 {
		condition["order_type"] = req.IsGrant
	}

	// 查询渠道订单列表
	list, total, err := l.svcCtx.ChannelOrderModel.FindPageByCondition(l.ctx, req.Page, req.Limit, keywordOrderNo, condition)
	if err != nil {
		logc.Errorf(l.ctx, "渠道订单关联渠道子订单查询失败，Err:%s", err)
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}

	if len(list) == 0 {
		return resp, nil
	}

	orderIds, _ := util.ArrayColumn(list, "Id")
	orderIds = util.ArrayUniqueValue(orderIds)
	// 获取订单号的所有子订单
	subOrderMap, err := l.svcCtx.ChannelSubOrderModel.BatchByChannelOrderIds(l.ctx, orderIds)
	if err != nil {
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}

	groupIds, _ := util.ArrayColumn(list, "BenefitsGroupId")
	groupIds = util.ArrayUniqueValue(groupIds)
	// 获取订单号的所有权益名称
	groupMap, err := l.svcCtx.BenefitGroupModel.BatchByGroupIds(l.ctx, groupIds)
	if err != nil {
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}

	// 发送短信的手机号列表
	var (
		phoneArr []any
		// phones   []string
	)
	phoneMap := make(map[string]string)
	for _, v := range list {
		if v.MaskPhone == "" {
			continue
		}

		// 解密手机号
		maskPhone, _ := util.AesDecrypt(l.ctx, v.MaskPhone, l.svcCtx.Config.EncryptKey)
		if _, exists := phoneMap[v.MaskPhone]; !exists {
			phoneMap[v.MaskPhone] = maskPhone
			phoneArr = append(phoneArr, maskPhone) // 需要收集到 phones 数组
		}
	}
	// if len(phoneArr) > 0 {
	// 	phones, _ = l.svcCtx.SysSmsRecordModel.FindByPhones(l.ctx, define.SmsLoginTemplateId, phoneArr)
	// 	if err != nil {
	// 		logc.Errorf(l.ctx, "短信发送记录查询失败，Err:%s", err)
	// 		return nil, errs.WithCode(err, errs.ServerErrorCode)
	// 	}
	// }

	for _, v := range list {
		info := types.OrderBenefitInfo{
			Id:        v.Id,
			IsGrant:   len(subOrderMap[v.Id]) > 0, // 大于0，代表已发放过权益
			OrderNo:   v.OrderNo,
			GroupName: groupMap[v.BenefitsGroupId].Name,
			Phone:     util.GenerateMaskPhone(phoneMap[v.MaskPhone]),
			UserType:  v.UserType,
			OrderType: v.OrderType,
			Source:    define.BenefitsSourceMapInt[v.Source],
			// IsSmsNotifiy: util.IsExist(phoneMap[v.MaskPhone], phones), // 是否发送过短信通知
			IsSmsNotifiy: true, // 是否发送过短信通知
			CreatedAt:    v.CreatedAt.Unix(),
			UpdatedAt:    v.UpdatedAt.Unix(),
		}
		resp.List = append(resp.List, info)
	}

	// 分页信息
	resp.Pagination = types.Pagination{
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	}
	return
}
