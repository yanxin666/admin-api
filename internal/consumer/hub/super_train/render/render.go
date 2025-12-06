package render

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/logz"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"muse-admin/internal/define"
	"muse-admin/internal/model/kingclub/supertrain"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
)

type Render struct {
	svcCtx *svc.ServiceContext
}

func NewRender(svcCtx *svc.ServiceContext) *Render {
	return &Render{
		svcCtx: svcCtx,
	}
}

func (r *Render) checkData(f *types.SuperTrainData) error {
	if f.CourseType < 0 {
		return errors.New("invalid course type")
	}

	// 目前只接收 1-2 的课程类型 1.写作 2.阅读
	if f.CourseType != 1 && f.CourseType != 2 {
		return errors.New("invalid course type")
	}

	// 章节老师的ID不能为空，会影响列表的展示
	for _, v := range f.Chapter {
		if v.TeacherId == 0 {
			return errors.New("invalid teacherId")
		}
	}

	return nil
}

// Render 处理数据
func (r *Render) Render(ctx context.Context, f *types.SuperTrainData) (err error) {
	// 数据校验
	if err = r.checkData(f); err != nil {
		return err
	}

	// 事务处理
	err = r.processTx(ctx, f)
	if err != nil {
		return err
	}

	// 将此条mq消费掉
	return nil
}

func (r *Render) processTx(ctx context.Context, f *types.SuperTrainData) error {
	// 开启事务：
	// 1. 插入课程表：ac_course 获取课程ID
	// 2. 插入章节表：ac_course_chapter 获取章节ID
	// 3. 插入任务表：ac_chapter_task 获取任务ID
	// 4. 插入子任务表：ac_sub_task 获取子任务ID
	// 5. 插入奖牌表：ac_medal 获取奖牌ID
	// 6. 插入奖牌与子任务的关联表：ac_sub_task_medal 获取关联ID
	err := r.svcCtx.MysqlConnKingClub.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		var (
			err      error
			courseId int64
		)

		// 1. 插入课节表：aw_course_outline 获取大纲ID
		c, err := convertCourse(r.svcCtx.Config.Mode, f)
		if err != nil {
			return err
		}
		courseId, err = r.handleCourse(ctx, session, c, f)
		if err != nil {
			return err
		}
		if courseId == 0 {
			// 数据过滤
			if err == nil {
				return nil
			}
			return errors.New("courseId为空")
		}

		// 2. 插入章节表：ac_course_chapter 获取章节ID
		// 一个课程中可能有多个章节
		for _, chapter := range f.Chapter {
			l, err := convertChapter(f.CourseType, chapter, courseId)
			if err != nil {
				return err
			}
			chapterId, err := r.handleChapter(ctx, session, l, f)
			if err != nil {
				return err
			}

			// 3. 插入任务表：ac_chapter_task 获取任务ID
			// 一个章节中可能有多个任务
			for _, task := range chapter.Task {
				t, err := convertTask(f.CourseType, task, courseId, chapterId)
				if err != nil {
					return err
				}
				taskId, err := r.handleTask(ctx, session, t, f)
				if err != nil {
					return err
				}

				// 4. 插入子任务表：ac_sub_task 获取子任务ID
				// 一个任务中可能有多个子任务
				for _, subTask := range task.SubTask {
					s, err := convertSubTask(f.CourseType, task.Extra.ParagraphType, subTask, taskId)
					if err != nil {
						return err
					}
					subtaskId, err := r.handleSubTask(ctx, session, s, f)
					if err != nil {
						return err
					}

					// 5. 插入奖牌表：ac_medal 获取奖牌ID
					if subTask.Medal.Name != "" || subTask.Medal.Description != "" || subTask.Medal.Type != 0 || subTask.Medal.PromptKnowledge != "" {
						m, err := convertMetal(subTask.Medal)
						if err != nil {
							return err
						}
						metalId, err := r.handleMetal(ctx, session, m)
						if err != nil {
							return err
						}

						// 6. 插入奖牌与子任务的关联表：ac_sub_task_medal 获取关联ID
						// 注意：插入之前需要先清空掉子任务的关联奖牌，目的是防止重复插入，因为每次推送会生成新奖牌
						err = r.delOldMedalRelation(ctx, session, subtaskId)
						if err != nil {
							return err
						}
						err = r.handleMedalRelation(ctx, session, subtaskId, metalId)
						if err != nil {
							return err
						}
					}
				}
			}
		}

		return nil
	})

	// 事务出错，记录日志
	if err != nil {
		return err
	}

	// ...
	return nil
}

func (r *Render) handleCourse(ctx context.Context, session sqlx.Session, c *supertrain.Course, f *types.SuperTrainData) (courseId int64, err error) {
	var data *supertrain.Course
	data, err = r.svcCtx.TrainCourseModel.FindOneByCourseNo(ctx, c.CourseNo)
	if err != nil {
		return 0, err
	}

	// 无课节数据
	if data == nil {
		// 传入状态为 下架，需要对其过滤
		if f.Status == define.ReviewStatus.OffShelf {
			logz.Infof(ctx, "要下架的课程课程编号: %s；课程名称: %s；数据不存在，需要过滤", f.CourseNo, f.CourseName)
			return 0, nil
		}
		courseId, err = r.svcCtx.TrainCourseModel.InsertSession(ctx, session, c)
		if err != nil {
			return 0, err
		}

		return courseId, nil
	}

	// 有课节数据，强制对其进行update
	c.Id = data.Id

	err = r.svcCtx.TrainCourseModel.UpdateSession(ctx, session, c)
	if err != nil {
		return 0, err
	}

	return data.Id, nil
}

func (r *Render) handleChapter(ctx context.Context, session sqlx.Session, l *supertrain.CourseChapter, f *types.SuperTrainData) (chapterId int64, err error) {
	// 可能存在同一个章节编号对应多个课程的情况
	data, err := r.svcCtx.TrainCourseChapterModel.FindOneByCourseIdAndChapterNo(ctx, l.CourseId, l.ChapterNo)
	if err != nil {
		return 0, err
	}

	// 无题目数据
	if data == nil {
		// 传入状态为 下架，需要对其过滤
		if f.Status == define.ReviewStatus.OffShelf {
			return 0, nil
		}
		chapterId, err = r.svcCtx.TrainCourseChapterModel.InsertSession(ctx, session, l)
		if err != nil {
			return 0, err
		}

		return chapterId, nil
	}

	// 有题目数据，强制对其进行update
	l.Id = data.Id

	err = r.svcCtx.TrainCourseChapterModel.UpdateSession(ctx, session, l)
	if err != nil {
		return 0, err
	}

	return data.Id, nil
}

func (r *Render) handleTask(ctx context.Context, session sqlx.Session, t *supertrain.ChapterTask, f *types.SuperTrainData) (taskId int64, err error) {
	// 可能存在同一个任务编号对应多个章节的情况
	data, err := r.svcCtx.TrainChapterTaskModel.FindOneByCourseIdAndChapterIdAndTaskNo(ctx, t.CourseId, t.CourseChapterId, t.TaskNo)
	if err != nil {
		return 0, err
	}

	// 无题目数据
	if data == nil {
		// 传入状态为 下架，需要对其过滤
		if f.Status == define.ReviewStatus.OffShelf {
			return 0, nil
		}
		taskId, err = r.svcCtx.TrainChapterTaskModel.InsertSession(ctx, session, t)
		if err != nil {
			return 0, err
		}

		return taskId, nil
	}

	// 有题目数据，强制对其进行update
	t.Id = data.Id

	err = r.svcCtx.TrainChapterTaskModel.UpdateSession(ctx, session, t)
	if err != nil {
		return 0, err
	}

	return data.Id, nil
}

func (r *Render) handleSubTask(ctx context.Context, session sqlx.Session, s *supertrain.SubTask, f *types.SuperTrainData) (subTaskId int64, err error) {
	data, err := r.svcCtx.TrainSubTaskModel.FindOneByTaskIdAndSubTaskNo(ctx, s.ChapterTaskId, s.SubTaskNo)
	if err != nil {
		return 0, err
	}

	// 无题目数据
	if data == nil {
		// 传入状态为 下架，需要对其过滤
		if f.Status == define.ReviewStatus.OffShelf {
			return 0, nil
		}
		subTaskId, err = r.svcCtx.TrainSubTaskModel.InsertSession(ctx, session, s)
		if err != nil {
			return 0, err
		}

		return subTaskId, nil
	}

	// 有题目数据，强制对其进行update
	s.Id = data.Id

	err = r.svcCtx.TrainSubTaskModel.UpdateSession(ctx, session, s)
	if err != nil {
		return 0, err
	}

	return data.Id, nil
}

func (r *Render) handleMetal(ctx context.Context, session sqlx.Session, s *supertrain.Medal) (medalId int64, err error) {
	medalId, err = r.svcCtx.TrainMedalModel.InsertSession(ctx, session, s)
	if err != nil {
		return 0, err
	}

	return medalId, nil
}

func (r *Render) delOldMedalRelation(ctx context.Context, session sqlx.Session, subtaskId int64) error {
	if subtaskId < 0 {
		return nil
	}

	err := r.svcCtx.TrainSubTaskMedalModel.DeleteSessionBySubTaskId(ctx, session, subtaskId)
	if err != nil {
		return err
	}

	return nil
}

func (r *Render) handleMedalRelation(ctx context.Context, session sqlx.Session, subtaskId, metalId int64) error {
	if metalId < 0 && subtaskId < 0 {
		return nil
	}

	_, err := r.svcCtx.TrainSubTaskMedalModel.InsertSession(ctx, session, &supertrain.SubTaskMedal{
		SubTaskId: subtaskId,
		MedalId:   metalId,
	})
	if err != nil {
		return err
	}

	return nil
}
