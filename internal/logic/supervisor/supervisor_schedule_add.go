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

type SupervisorScheduleAdd struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSupervisorScheduleAdd(ctx context.Context, svcCtx *svc.ServiceContext) *SupervisorScheduleAdd {
	return &SupervisorScheduleAdd{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SupervisorScheduleAdd) SupervisorScheduleAdd(req *types.SupervisorScheduleEditReq) (resp *types.Result, err error) {
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
	// 1. 新增直播表，获取直播自增ID
	// 2. 新增督学排课表
	// 3. 互动关联表 2025.11.20此逻辑废弃
	err = l.svcCtx.MysqlConnKingClub.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		var (
			streamId int64
			// scheduleId int64
		)

		// 1. 新增直播表，获取直播自增ID：lv_stream
		streamId, err = l.svcCtx.SupeStreamModel.InsertSession(ctx, session, &supervisor.Stream{
			Title:     req.StreamTitle,
			StartTime: readyTime,
		})
		if err != nil {
			return err
		}

		// 2. 新增排课表：supe_schedule
		schedule := l.convertSchedule(req, streamId, effectiveStart, effectiveEnd, startTime, endTime, readyTime)
		_, err = l.svcCtx.SupeScheduleModel.InsertSession(ctx, session, schedule)
		if err != nil {
			return err
		}

		// // 3. 插入督学互动关联表：supe_schedule_interaction_link
		// streamId, err = l.svcCtx.SupeScheduleInteractRelationModel.InsertSession(ctx, session, &supervisor.ScheduleInteractionLink{
		// 	ScheduleId:    scheduleId,
		// 	InteractionId: req.InteractId,
		// 	Status:        1,
		// })
		// if err != nil {
		// 	return err
		// }

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

func checkDate(effectiveStartData, effectiveEndData, startDate, endDate, readyDate string) (effectiveStart, effectiveEnd, startTime, endTime, readyTime time.Time, err error) {
	effectiveStart, err = util.GetConverterDate(effectiveStartData, util.StandardDatetime)
	if err != nil {
		return effectiveStart, effectiveEnd, startTime, endTime, readyTime, errs.NewMsg(errs.ErrCodeParamsAbnormal, "开始时间日期格式错误").ShowMsg()
	}

	effectiveEnd, err = util.GetConverterDate(effectiveEndData, util.StandardDatetime)
	if err != nil {
		return effectiveStart, effectiveEnd, startTime, endTime, readyTime, errs.NewMsg(errs.ErrCodeParamsAbnormal, "开始时间日期格式错误").ShowMsg()
	}

	startTime, err = util.GetConverterDate(startDate, util.StandardDatetime)
	if err != nil {
		return effectiveStart, effectiveEnd, startTime, endTime, readyTime, errs.NewMsg(errs.ErrCodeParamsAbnormal, "开始时间日期格式错误").ShowMsg()
	}

	endTime, err = util.GetConverterDate(endDate, util.StandardDatetime)
	if err != nil {
		return effectiveStart, effectiveEnd, startTime, endTime, readyTime, errs.NewMsg(errs.ErrCodeParamsAbnormal, "结束时间日期格式错误").ShowMsg()
	}

	readyTime, err = util.GetConverterDate(readyDate, util.StandardDatetime)
	if err != nil {
		return effectiveStart, effectiveEnd, startTime, endTime, readyTime, errs.NewMsg(errs.ErrCodeParamsAbnormal, "结束时间日期格式错误").ShowMsg()
	}

	return effectiveStart, effectiveEnd, startTime, endTime, readyTime, nil
}

// 转换为supervisor_schedule表可录入的数据
func (l *SupervisorScheduleAdd) convertSchedule(req *types.SupervisorScheduleEditReq, streamId int64, effectiveStart, effectiveEnd, startTime, endTime, readyTime time.Time) *supervisor.Schedule {
	return &supervisor.Schedule{
		TeacherId:        req.TearcherId,
		TeacherName:      req.TeacherName,
		CourseId:         req.CourseId,
		CourseType:       req.CourseType,
		Name:             req.Name,
		StreamId:         streamId,
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
