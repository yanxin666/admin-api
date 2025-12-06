package session

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"muse-admin/pkg/errs"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetListByUserId struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetListByUserId(ctx context.Context, svcCtx *svc.ServiceContext) *GetListByUserId {
	return &GetListByUserId{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetListByUserId  获取用户的会话列表
func (l *GetListByUserId) GetListByUserId(req *types.GetListByUserIdReq) (resp *types.GetListByUserIdResp, err error) {
	resp = &types.GetListByUserIdResp{}
	// 校验参数
	if req.UserId <= 0 {
		return nil, errs.NewMsg(errs.ErrCodeProgram, "参数错误")
	}
	// 查询数据
	res, err := l.svcCtx.SessionModel.FindListByUserId(l.ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	// 暂无数据
	if len(res) == 0 {
		return nil, nil
	}
	// 构造数据
	list := make([]*types.EventSession, 0)
	for _, v := range res {
		// fmt.Println(fmt.Sprintf("%+v", v))
		list = append(list, &types.EventSession{
			StratTime: util.ConvertTimeToFormattedDate(v.StartTime),
			EndTime:   util.ConvertTimeToFormattedDate(v.EndTime),
			Session:   v.Session,
			SessionId: v.Id,
		})
	}
	// 构造返回结果
	return &types.GetListByUserIdResp{
		List: list,
	}, nil

}
