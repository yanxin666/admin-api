package course

import (
	"context"
	"muse-admin/internal/svc/course"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCourseChapter struct {
	logx.Logger
	ctx            context.Context
	svcCtx         *svc.ServiceContext
	chapterService *course.ChapterService
}

func NewGetCourseChapter(ctx context.Context, svcCtx *svc.ServiceContext) *GetCourseChapter {
	return &GetCourseChapter{
		Logger:         logx.WithContext(ctx),
		ctx:            ctx,
		svcCtx:         svcCtx,
		chapterService: course.NewChapterService(svcCtx),
	}
}

func (l *GetCourseChapter) GetCourseChapter(req *types.GetCourseChapterReq) (resp *types.GetCourseChapterResp, err error) {
	list, err := l.chapterService.GetCourseIdByChapterList(l.ctx, req.CourseId)
	if err != nil {
		return nil, err
	}
	resp = &types.GetCourseChapterResp{ChapterList: make([]types.CourseChapter, 0)}
	for _, chapter := range list {
		resp.ChapterList = append(resp.ChapterList, types.CourseChapter{
			Id:        chapter.Id,
			CourseId:  chapter.CourseId,
			Title:     chapter.Title,
			Name:      chapter.Name,
			ChapterNo: chapter.ChapterNo,
		})
	}
	return
}

// 校验参数
func (l *GetCourseChapter) checkParams(req *types.GetCourseChapterReq) error {
	// todo: add your logic check params and delete this line
	// 若校验的参数过少，可删除此方法，直接在上方 GetCourseChapter 中编写

	return nil
}
