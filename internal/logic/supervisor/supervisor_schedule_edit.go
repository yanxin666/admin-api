package supervisor

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"muse-admin/internal/model/supervisor"
	"muse-admin/pkg/errs"
	"time"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SupervisorScheduleEdit struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSupervisorScheduleEdit(ctx context.Context, svcCtx *svc.ServiceContext) *SupervisorScheduleEdit {
	return &SupervisorScheduleEdit{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SupervisorScheduleEdit) SupervisorScheduleEdit(req *types.SupervisorScheduleEditReq) (resp *types.Result, err error) {
	// 预约人数检查
	if req.SurplusStock > req.MaxStock {
		return nil, errs.NewMsg(errs.ErrCodeParamsAbnormal, "剩余预约人数不可大于预约总人数").ShowMsg()
	}

	// 时间格式检查
	effectiveStart, effectiveEnd, startTime, endTime, readyTime, err := checkDate(req.EffectiveTime[0], req.EffectiveTime[1], req.AppointmentTime[0], req.AppointmentTime[1], req.StartTime)
	if err != nil {
		return nil, err
	}

	// BgUrl和Content格式检查
	if !util.IsJSON(req.BgUrl) {
		return nil, errs.NewMsg(errs.ErrCodeParamsAbnormal, "背景数据格式不正确").ShowMsg()
	}
	if !util.IsJSON(req.Content) {
		return nil, errs.NewMsg(errs.ErrCodeParamsAbnormal, "课程内容格式不正确").ShowMsg()
	}

	// 事务操作：
	// 1. 修改直播表
	// 2. 修改督学排课表
	err = l.svcCtx.MysqlConnKingClub.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		var (
			scheduleData *supervisor.Schedule
			streamData   *supervisor.Stream
		)

		// 1. 修改排课表：supe_schedule
		scheduleData, err = l.svcCtx.SupeScheduleModel.FindOne(ctx, req.Id)
		if err != nil {
			return err
		}
		schedule := l.convertSchedule(scheduleData, req, effectiveStart, effectiveEnd, startTime, endTime, readyTime)
		err = l.svcCtx.SupeScheduleModel.UpdateSession(ctx, session, schedule)
		if err != nil {
			return err
		}

		// 2. 修改直播表名称
		streamData, err = l.svcCtx.SupeStreamModel.FindOne(ctx, scheduleData.StreamId)
		if err != nil {
			return err
		}
		streamData.Title = req.StreamTitle
		err = l.svcCtx.SupeStreamModel.UpdateSession(ctx, session, streamData)
		if err != nil {
			return err
		}

		return nil
	})

	// 事务出错，记录日志
	if err != nil {
		return nil, err
	}

	return &types.Result{
		Result: true,
	}, nil
}

// 转换为supervisor_schedule表可录入的数据
func (l *SupervisorScheduleEdit) convertSchedule(data *supervisor.Schedule, req *types.SupervisorScheduleEditReq, effectiveStart, effectiveEnd, startTime, endTime, readyTime time.Time) *supervisor.Schedule {
	return &supervisor.Schedule{
		Id:               data.Id,
		TeacherId:        req.TearcherId,
		TeacherName:      req.TeacherName,
		CourseId:         req.CourseId,
		CourseType:       req.CourseType,
		Name:             req.Name,
		StreamId:         data.StreamId,
		EffectiveStart:   effectiveStart,
		EffectiveEnd:     effectiveEnd,
		AppointmentStart: startTime,
		AppointmentEnd:   endTime,
		StartTime:        readyTime,
		MaxStock:         req.MaxStock,
		SurplusStock:     req.SurplusStock,
		ScheduleType:     req.ScheduleType,
		Status:           req.Status,
		BgUrl:            req.BgUrl,
		Content:          req.Content,
	}
}
