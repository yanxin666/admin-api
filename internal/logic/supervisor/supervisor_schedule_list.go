package supervisor

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"github.com/zeromicro/go-zero/core/logc"
	"muse-admin/internal/model/tools"
	"muse-admin/pkg/errs"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SupervisorScheduleList struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSupervisorScheduleList(ctx context.Context, svcCtx *svc.ServiceContext) *SupervisorScheduleList {
	return &SupervisorScheduleList{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SupervisorScheduleList) SupervisorScheduleList(req *types.SupervisorListReq) (resp *types.SupervisorListResp, err error) {
	resp = &types.SupervisorListResp{
		List:       make([]types.SupervisorInfo, 0),
		Pagination: types.Pagination{},
	}
	// 构造查询条件
	condition := tools.FilterConditions(req)
	// 获取督学排课列表
	list, total, err := l.svcCtx.SupeScheduleModel.FindPageByCondition(l.ctx, req.Page, req.Limit, condition)
	if err != nil {
		logc.Errorf(l.ctx, "批量查询督学排课列表失败，Err:%s", err)
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}
	// 空数据
	if len(list) == 0 {
		return resp, nil
	}
	// 批量获取课程信息
	courseIds, _ := util.ArrayColumn(list, "CourseId")
	courseInfos, err := l.svcCtx.TrainCourseModel.BatchMapByIds(l.ctx, courseIds)
	if err != nil {
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}

	// 批量获取直播表信息
	streamIds, _ := util.ArrayColumn(list, "StreamId")
	streamInfos, err := l.svcCtx.SupeStreamModel.BatchMapByIds(l.ctx, streamIds)
	if err != nil {
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}
	// 拼装返回数据
	for _, item := range list {
		resp.List = append(resp.List, types.SupervisorInfo{
			Id:               item.Id,
			TeacherId:        item.TeacherId,                   // 注: 教师ID 与TeacherName下面不对应
			TeacherName:      item.TeacherName,                 // 教师名字
			CourseId:         item.CourseId,                    // 课程ID
			CourseName:       courseInfos[item.CourseId].Name,  // 课程名字
			CourseType:       item.CourseType,                  // 课程类型 1.超练作文慢练 2.超练作文快练 3.超练阅读
			Name:             item.Name,                        // 督学排课名称
			StreamTitle:      streamInfos[item.StreamId].Title, // 直播标题
			ScheduleType:     item.ScheduleType,                // (排课类型:1.直播2.录播)
			StartTime:        item.StartTime.Unix(),
			AppointmentStart: item.AppointmentStart.Unix(),
			AppointmentEnd:   item.AppointmentEnd.Unix(),
			Status:           item.Status,       // (状态:1.待开放2.开放3.进行中4.结束)
			MaxStock:         item.MaxStock,     // 最大预约量
			SurplusStock:     item.SurplusStock, // 剩余预约量
			BgUrl:            item.BgUrl,        // 伴学、督学课程背景图片
			Content:          item.Content,      // 伴学、督学课程文案
			CreatedAt:        item.CreatedAt.Unix(),
			UpdatedAt:        item.UpdatedAt.Unix(),
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
