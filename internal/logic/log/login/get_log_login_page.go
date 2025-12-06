package login

import (
	"context"
	"muse-admin/internal/define"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetLogLoginPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLogLoginPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLogLoginPageLogic {
	return &GetLogLoginPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLogLoginPageLogic) GetLogLoginPage(req *types.LogLoginPageReq) (resp *types.LogLoginPageResp, err error) {
	loginLogList, err := l.svcCtx.SysLogModel.FindPage(l.ctx, define.SysLoginLogType, req.Page, req.Limit)
	if err != nil {
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}

	var loginLog types.LogLogin
	logList := make([]types.LogLogin, 0)
	for _, v := range loginLogList {
		err := copier.Copy(&loginLog, &v)
		loginLog.CreateTime = v.CreateTime.Format(define.SysDateFormat)
		if err != nil {
			return nil, errs.WithCode(err, errs.ServerErrorCode)
		}
		logList = append(logList, loginLog)
	}

	total, err := l.svcCtx.SysLogModel.FindCount(l.ctx, define.SysLoginLogType)
	if err != nil {
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}

	pagination := types.Pagination{
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	}

	return &types.LogLoginPageResp{
		List:       logList,
		Pagination: pagination,
	}, nil
}
