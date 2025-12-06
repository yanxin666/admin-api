package dict

import (
	"context"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetConfigDictPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetConfigDictPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConfigDictPageLogic {
	return &GetConfigDictPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetConfigDictPageLogic) GetConfigDictPage(req *types.ConfigDictPageReq) (resp *types.ConfigDictPageResp, err error) {
	configDictionaryList, err := l.svcCtx.SysDictionaryModel.FindPageByParentId(l.ctx, req.ParentId, req.Page, req.Limit)
	if err != nil {
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}

	var dictionary types.ConfigDict
	dictionaryList := make([]types.ConfigDict, 0)
	for _, sysDictionary := range configDictionaryList {
		err := copier.Copy(&dictionary, &sysDictionary)
		if err != nil {
			return nil, errs.WithCode(err, errs.ServerErrorCode)
		}
		dictionaryList = append(dictionaryList, dictionary)
	}

	total, err := l.svcCtx.SysDictionaryModel.FindCountByParentId(l.ctx, req.ParentId)
	if err != nil {
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}

	pagination := types.Pagination{
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	}

	return &types.ConfigDictPageResp{
		List:       dictionaryList,
		Pagination: pagination,
	}, nil
}
