package course

import (
	"context"
	"muse-admin/internal/svc/course"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChapterTask struct {
	logx.Logger
	ctx            context.Context
	svcCtx         *svc.ServiceContext
	chapterService *course.ChapterService
}

func NewGetChapterTask(ctx context.Context, svcCtx *svc.ServiceContext) *GetChapterTask {
	return &GetChapterTask{
		Logger:         logx.WithContext(ctx),
		ctx:            ctx,
		svcCtx:         svcCtx,
		chapterService: course.NewChapterService(svcCtx),
	}
}

func (l *GetChapterTask) GetChapterTask(req *types.GetChapterTaskReq) (resp *types.GetChapterTaskResp, err error) {
	list, err := l.chapterService.GetChapterIdByTaskList(l.ctx, req.ChapterId)
	if err != nil {
		return nil, err
	}
	resp = &types.GetChapterTaskResp{TaskList: make([]types.ChapterTask, 0)}
	for _, task := range list {
		resp.TaskList = append(resp.TaskList, types.ChapterTask{
			Id:        task.Id,
			ChapterId: task.CourseChapterId,
			Title:     task.Title,
			TaskNo:    task.TaskNo,
		})
	}

	return
}
