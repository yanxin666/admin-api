package supervisor

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"muse-admin/internal/define"
	"muse-admin/internal/model/supervisor"
	"muse-admin/pkg/errs"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SupervisorScheduleStatus struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSupervisorScheduleStatus(ctx context.Context, svcCtx *svc.ServiceContext) *SupervisorScheduleStatus {
	return &SupervisorScheduleStatus{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SupervisorScheduleStatus) SupervisorScheduleStatus(req *types.SupervisorScheduleStatusReq) error {
	// 检查参数
	if !util.InSlice(req.Status, []int64{
		define.KCScheduleCourseStatus.Pending,
		define.KCScheduleCourseStatus.Ongoing,
		define.KCScheduleCourseStatus.Open,
		define.KCScheduleCourseStatus.Finished,
	}) {
		return errs.NewMsg(errs.WarningCode, "状态参数不合法").ShowMsg()
	}
	// 更新督学排课状态
	_, err := l.svcCtx.SupeScheduleModel.UpdateFillFieldsById(l.ctx, req.Id, &supervisor.Schedule{
		Status: req.Status,
	})
	if err != nil {
		return errs.WithMsg(err, errs.ErrCodeAbnormal, "更新督学排课状态失败")
	}
	return nil
}
