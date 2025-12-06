package supervisor

import (
	"context"
	"muse-admin/internal/define"
	"muse-admin/internal/model/tools"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseList struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseList(ctx context.Context, svcCtx *svc.ServiceContext) *CourseList {
	return &CourseList{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CourseList) CourseList(req *types.CourseListReq) (resp *types.CourseListResp, err error) {
	resp = &types.CourseListResp{}
	// 构造查询条件
	condition := tools.FilterConditions(req)
	// 获取课程列表
	list, err := l.svcCtx.TrainCourseModel.FindListByConds(l.ctx, define.ReviewStatus.Passed, condition)
	if err != nil {
		return nil, err
	}
	resp.List = make([]types.Course, 0)
	for _, v := range list {
		resp.List = append(resp.List, types.Course{
			Id:         v.Id,
			CourseType: v.CourseType, // 课程类型 1.超练作文慢练 2.超练作文快练 3.超练阅读
			Name:       v.Name,
			LessonName: v.LessonName,
		})
	}
	return
}
