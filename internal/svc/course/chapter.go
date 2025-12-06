package course

import (
	"context"
	"muse-admin/internal/model/kingclub/supertrain"
	"muse-admin/internal/svc"
)

type ChapterService struct {
	svcCtx *svc.ServiceContext
}

func NewChapterService(svcCtx *svc.ServiceContext) *ChapterService {
	return &ChapterService{
		svcCtx: svcCtx,
	}
}

// GetCourseIdByChapterList 根据courseID获取章节列表
func (c *ChapterService) GetCourseIdByChapterList(ctx context.Context, courseId int64) (resp []supertrain.CourseChapter, err error) {
	return c.svcCtx.TrainCourseChapterModel.GetCourseIdByChapterList(ctx, courseId)
}

// GetChapterIdByTaskList 根据章节Id获取任务列表
func (c *ChapterService) GetChapterIdByTaskList(ctx context.Context, chapterId int64) (resp []supertrain.ChapterTask, err error) {
	return c.svcCtx.TrainChapterTaskModel.GetChapterIdByTaskList(ctx, chapterId)
}
