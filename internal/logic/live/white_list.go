package live

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"github.com/zeromicro/go-zero/core/logc"
	"muse-admin/pkg/errs"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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

func (l *WhiteList) WhiteList(req *types.LiveBetaUserListReq) (resp *types.LiveBetaUserListResp, err error) {
	var (
		phoneKeyword, usernameKeyword, companyKeyword string
	)

	resp = &types.LiveBetaUserListResp{
		List:       make([]types.LiveBetaUserInfo, 0),
		Pagination: types.Pagination{},
	}

	// 构建查询条件
	condition := make(map[string]interface{})
	if req.Id != 0 {
		condition["id"] = req.Id
	}
	if req.Status != 0 {
		condition["status"] = req.Status
	}
	if len(req.Phone) != 0 {
		phoneKeyword = req.Phone
	}
	if len(req.Username) != 0 {
		usernameKeyword = req.Username
	}
	if len(req.Company) != 0 {
		companyKeyword = req.Company
	}

	// 获取申请用户列表
	list, total, err := l.svcCtx.LiveBetaRecordModel.FindPageByCondition(l.ctx, req.Page, req.Limit, phoneKeyword, usernameKeyword, companyKeyword, condition)
	if err != nil {
		logc.Errorf(l.ctx, "批量查询用户失败，Err:%s", err)
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}

	// 空数据
	if len(list) == 0 {
		return resp, nil
	}

	// 获取操作人名称
	operateIds, _ := util.ArrayColumn(list, "OperateId")
	operateIds = util.ArrayUniqueValue(operateIds)
	userInfo, err := l.svcCtx.SysUserModel.BatchByUserIds(l.ctx, operateIds)
	if err != nil {
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}

	for _, v := range list {
		resp.List = append(resp.List, types.LiveBetaUserInfo{
			Id:          v.Id,
			Phone:       v.Phone,
			Username:    v.UserName,
			Company:     v.Company,
			Status:      v.Status,
			Remark:      v.Remark,
			OperateName: userInfo[v.OperateId].Username,
			OperateAt:   v.OperateAt,
			CreatedAt:   v.CreatedAt.Unix(),
			UpdatedAt:   v.UpdatedAt.Unix(),
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
