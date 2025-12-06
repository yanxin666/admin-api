package course

import (
	"context"

	"muse-admin/internal/define"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCourseCommonData struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCourseCommonData(ctx context.Context, svcCtx *svc.ServiceContext) *GetCourseCommonData {
	return &GetCourseCommonData{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCourseCommonData) GetCourseCommonData() (resp *types.CourseCommonData, err error) {
	// 转换系列类型列表
	seriesTypeList := make([]types.SeriesType, 0, len(define.KCSeriesTypeList))
	for _, item := range define.KCSeriesTypeList {
		seriesTypeList = append(seriesTypeList, types.SeriesType{
			SeriesType:     item["series_type"].(int),
			SeriesTypeName: item["series_type_name"].(string),
		})
	}

	// 转换课程类型列表
	courseTypeList := make([]types.CourseType, 0, len(define.KCCourseTypeList))
	for _, item := range define.KCCourseTypeList {
		courseTypeList = append(courseTypeList, types.CourseType{
			CourseType:     item["course_type"].(int),
			CourseTypeName: item["course_type_name"].(string),
		})
	}

	// 转换任务类型列表
	taskTypeList := make([]types.TaskType, 0, len(define.KCTaskTypeList))
	for _, item := range define.KCTaskTypeList {
		taskTypeList = append(taskTypeList, types.TaskType{
			TaskType:     item["task_type"].(int),
			TaskTypeName: item["task_type_name"].(string),
		})
	}

	// 转换子任务类型列表
	subTaskTypeList := make([]types.SubTaskType, 0, len(define.KCSubTaskTypeList))
	for _, item := range define.KCSubTaskTypeList {
		subTaskTypeList = append(subTaskTypeList, types.SubTaskType{
			SubTaskType:     item["sub_task_type"].(int),
			SubTaskTypeName: item["sub_task_type_name"].(string),
		})
	}

	// 转换教师ID列表
	teacherIdList := make([]types.Teacher, 0, len(define.KCTeacherIdList))
	for _, item := range define.KCTeacherIdList {
		teacherIdList = append(teacherIdList, types.Teacher{
			Id:   item["id"].(int),
			Name: item["name"].(string),
		})
	}

	resp = &types.CourseCommonData{
		SeriesTypeList:  seriesTypeList,
		CourseTypeList:  courseTypeList,
		TeacherIds:      teacherIdList,
		TaskTypeList:    taskTypeList,
		SubTaskTypeList: subTaskTypeList,
	}
	return
}
