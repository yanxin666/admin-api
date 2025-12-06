package member

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"muse-admin/internal/model/tools"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"
)

type WhiteList struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWhiteList(ctx context.Context, svcCtx *svc.ServiceContext) *WhiteList {
	return &WhiteList{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WhiteList) WhiteList(req *types.WhiteListReq) (resp *types.WhiteListResp, err error) {
	var (
		phone string
	)
	resp = &types.WhiteListResp{
		List:       make([]types.WhiteInfo, 0),
		Pagination: types.Pagination{},
	}

	// 构建查询条件
	condition := tools.FilterConditions(req)
	// 删除不需要的查询条件
	util.MapDeleteExistKeys(condition, []string{"phone", "start_time", "end_time", "remark"})
	if len(req.Phone) == 11 {
		condition["phone"] = req.Phone
	} else {
		phone = req.Phone
	}

	// 获取主用户列表
	list, total, err := l.svcCtx.UserWhiteModel.FindPageByCondition(l.ctx, req.Page, req.Limit, phone, req.StartTime, req.EndTime, req.Remark, condition)
	if err != nil {
		logc.Errorf(l.ctx, "批量查询用户失败，Err:%s", err)
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}

	// 空数据
	if len(list) == 0 {
		return resp, nil
	}

	for _, v := range list {
		resp.List = append(resp.List, types.WhiteInfo{
			Id:        v.Id,
			Phone:     v.Phone,
			StartTime: v.StartTime.Unix(),
			EndTime:   v.EndTime.Unix(),
			Product:   v.Product,
			Source:    v.Source,
			Status:    v.Status,
			Remark:    v.Remark,
			CreatedAt: v.CreatedAt.Unix(),
			UpdatedAt: v.UpdatedAt.Unix(),
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
