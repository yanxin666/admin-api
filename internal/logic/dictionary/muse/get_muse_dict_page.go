package muse

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"muse-admin/internal/model/tools"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMuseDictPage struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMuseDictPage(ctx context.Context, svcCtx *svc.ServiceContext) *GetMuseDictPage {
	return &GetMuseDictPage{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMuseDictPage) GetMuseDictPage(req *types.MuseDictPageReq) (resp *types.MuseDictPageResp, err error) {
	resp = &types.MuseDictPageResp{
		List:       make([]types.SysDictWhite, 0),
		Pagination: types.Pagination{},
	}

	// 构建查询条件
	condition := tools.FilterConditions(req)
	// 删除不需要的查询条件
	util.MapDeleteExistKeys(condition, []string{"key", "remark"})

	res, total, err := l.svcCtx.DictWhiteModel.FindPageByCondition(l.ctx, req.Page, req.Limit, req.Key, req.Remark, condition)
	if err != nil {
		return nil, err
	}

	for _, item := range res {
		resp.List = append(resp.List, types.SysDictWhite{
			Id:        item.Id,
			Key:       item.Key,
			Remark:    item.Remark,
			CreatedAt: item.CreatedAt.Unix(),
			UpdatedAt: item.UpdatedAt.Unix(),
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
