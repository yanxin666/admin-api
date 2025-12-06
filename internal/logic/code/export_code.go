package code

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/logz"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	ctxt "muse-admin/internal/tools"
	"muse-admin/pkg/excel"
	"strings"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExportCode struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExportCode(ctx context.Context, svcCtx *svc.ServiceContext) *ExportCode {
	return &ExportCode{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExportCode) ExportCode(req *types.ListCodeReq) (resp *types.ExportCodeResp, err error) {
	if err = l.checkParams(req); err != nil {
		return nil, err
	}
	condition := l.ExportConditionParams(req)
	page := int64(1)
	pageSize := int64(10000)
	if req.Page != 0 {
		page = req.Page
	}
	if req.PageSize != 0 {
		pageSize = req.PageSize
	}

	generateCodeData, _, err := l.svcCtx.RedeemCodeModel.FindPageByCondition(l.ctx, page, pageSize, condition)
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

	// 设置表头的映射
	headerMap := map[string]string{
		"Code":            "兑换码",
		"BenefitsGroupId": "权益组ID",
		"ValidDate":       "过期时间",
		"Batch":           "批次",
		"Status":          "状态:1.未兑换 2.兑换",
	}

	o := excel.NewExCelExport(
		l.svcCtx.Config.Oss,
		excel.SetFilepath("code"),
		excel.SetFilename(util.GetNowTimeNoFormat()+"-"+ctxt.GetUserInfoByCtx(l.ctx).Username),
	)

	url, err := o.GenExcelFile(l.ctx, headerMap, list)
	if err != nil {
		logz.Errorf(l.ctx, "[code列表导出]upload to cos failed, err:%v", err)
		return nil, err
	}

	newUrl := strings.Replace(url, "https://static-1318590712.cos.ap-beijing.myqcloud.com", "https://static.zmexing.com", 1)

	return &types.ExportCodeResp{
		Url: newUrl,
	}, nil
}

func (l *ExportCode) ExportConditionParams(req *types.ListCodeReq) map[string]interface{} {
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
func (l *ExportCode) checkParams(req *types.ListCodeReq) error {
	// todo: add your logic check params and delete this line
	// 若校验的参数过少，可删除此方法，直接在上方 ExportCode 中编写

	return nil
}
