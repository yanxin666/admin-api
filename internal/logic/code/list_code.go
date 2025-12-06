package code

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCode struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListCode(ctx context.Context, svcCtx *svc.ServiceContext) *ListCode {
	return &ListCode{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListCode) ListCode(req *types.ListCodeReq) (resp *types.ListCodeResp, err error) {
	if err = l.checkParams(req); err != nil {
		return nil, err
	}
	condition := l.conditionParams(req)
	page := int64(1)
	pageSize := int64(20)
	if req.Page != 0 {
		page = req.Page
	}
	if req.PageSize != 0 {
		pageSize = req.PageSize
	}
	generateCodeData, total, err := l.svcCtx.RedeemCodeModel.FindPageByCondition(l.ctx, page, pageSize, condition)
	if err != nil {
		return nil, err
	}

	var list []types.CodeList
	for _, v := range generateCodeData {
		validDate := util.ConvertTimeToFormattedDate(v.ValidDate)
		codeInfo := types.CodeList{
			Id:              v.Id,
			Code:            v.Code,
			BenefitsGroupId: v.BenefitsGroupId,
			ValidDate:       validDate,
			Batch:           v.Batch,
			Status:          v.Status,
		}
		list = append(list, codeInfo)
	}
	return &types.ListCodeResp{
		List: list,
		Pagination: types.ListCodePagination{
			Page:  page,
			Limit: pageSize,
			Total: total,
		},
	}, nil

}

func (l *ListCode) conditionParams(req *types.ListCodeReq) map[string]interface{} {
	condition := make(map[string]interface{})
	if req.BenefitsGroupId != 0 {
		condition["benefits_group_id"] = req.BenefitsGroupId
	}
	if req.Batch != 0 {
		condition["batch"] = req.Batch
	}

	if req.Code != "" {
		condition["code"] = req.Code
	}

	return condition
}

// 校验参数
func (l *ListCode) checkParams(req *types.ListCodeReq) error {
	// todo: add your logic check params and delete this line
	// 若校验的参数过少，可删除此方法，直接在上方 ListCode 中编写

	return nil
}
