package member

import (
	"context"
	userModel "muse-admin/internal/model/user"
	"muse-admin/internal/svc"
	ctxt "muse-admin/internal/tools"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type WhiteEdit struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWhiteEdit(ctx context.Context, svcCtx *svc.ServiceContext) *WhiteEdit {
	return &WhiteEdit{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WhiteEdit) WhiteEdit(req *types.WhiteEditReq) (resp *types.Result, err error) {
	if err = l.checkParams(req); err != nil {
		return nil, err
	}

	// 获取手机号 以及 手机号姓名的映射
	phones, zhCnMap, err := getPhones(req.Phone)
	if err != nil {
		return nil, err
	}
	if len(phones) != 1 {
		return nil, errs.NewMsg(errs.ErrCodeParamsAbnormal, "编辑时手机号只能为单条！").ShowMsg()
	}

	// 时间格式检查
	startTime, endTime, err := checkDate(req.Date[0], req.Date[1])
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.UserWhiteModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, err
	}

	remark := "【" + ctxt.GetUserInfoByCtx(l.ctx).Username + "】"
	if req.Remark != "" {
		remark = remark + req.Remark
	}

	// 如果手机号有对应的姓名，则取出
	if zhCnMap != nil && zhCnMap[phones[0]] != "" {
		remark = remark + "-" + zhCnMap[phones[0]]
	}

	err = l.svcCtx.UserWhiteModel.Update(l.ctx, &userModel.WhiteList{
		Id:        req.Id,
		Phone:     phones[0],
		StartTime: startTime,
		EndTime:   endTime,
		Product:   req.Product,
		Source:    req.Source,
		Status:    req.Status,
		Remark:    remark,
	})
	if err != nil {
		return nil, err
	}

	return &types.Result{
		Result: true,
	}, nil
}

// 校验参数
func (l *WhiteEdit) checkParams(req *types.WhiteEditReq) error {
	if req.Id <= 0 {
		return errs.NewMsg(errs.ErrCodeParamsAbnormal, "ID不存在").ShowMsg()
	}

	return nil
}
