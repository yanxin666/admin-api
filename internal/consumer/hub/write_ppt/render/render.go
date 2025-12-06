package render

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"muse-admin/internal/define"
	"muse-admin/internal/model/write_ppt"
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

// Render 处理数据
func (r *Render) Render(ctx context.Context, f *types.WritePPTData) (err error) {
	// todo 不接受下架数据 原因：技巧赏析是原子操作，无法做到统一改动，需要具体细节商订
	if f.ReviewStatus == 5 {
		return nil
	}

	if len(f.PPTFiles) == 0 {
		return errors.New("ppt_files字段为空")
	}

	// 事务处理
	err = r.processTx(ctx, f)
	if err != nil {
		return err
	}

	// 将此条mq消费掉
	return nil
}

func (r *Render) processTx(ctx context.Context, f *types.WritePPTData) error {
	// 开启事务：
	// 1. 插入课节表：aw_course_outline 获取大纲ID
	// 2. 插入课节表：aw_lesson
	err := r.svcCtx.MysqlConnAbility.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		var (
			err      error
			courseId int64
		)

		// 1. 插入课节表：aw_course_outline 获取大纲ID
		c := convertCourse(f)
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

		// 2. 插入课节表：aw_lesson
		l := convertLesson(f, courseId)
		_, err = r.handleLesson(ctx, session, l, f)
		if err != nil {
			return err
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

func (r *Render) handleCourse(ctx context.Context, session sqlx.Session, c *write_ppt.CourseOutline, f *types.WritePPTData) (courseId int64, err error) {
	var data *write_ppt.CourseOutline
	data, err = r.svcCtx.PPTCourseOutlineModel.FindOneByUnitAndSeries(ctx, c.Unit, c.Series)
	if err != nil {
		return 0, err
	}

	// 因数据源为赏析的topic不准确，无法进行投产使用，所以需要手动处理
	if f.LessonType == 2 {
		c.Topic = ""
	}

	// 无课节数据
	if data == nil {
		// 传入状态为 下架，需要对其过滤
		if f.ReviewStatus == define.ReviewStatus.OffShelf {
			return 0, nil
		}
		// 传入状态为 审核通过，需要新增
		courseId, err = r.svcCtx.PPTCourseOutlineModel.InsertSession(ctx, session, c)
		if err != nil {
			return 0, err
		}

		return courseId, nil
	}

	// 有课节数据，强制对其进行update
	c.Id = data.Id
	c.Name = data.Name // 保留脚本生成的数据
	// 若为已上线的数据，需要保留其状态
	if data.Status == 2 {
		c.Status = 2
	}

	// 再次进来为技巧，获取技巧的topic
	if f.LessonType == 1 {
		c.Topic = f.Topic
	}
	// 再次进来为赏析，保留库中的topic
	if f.LessonType == 2 {
		c.Topic = data.Topic
	}

	err = r.svcCtx.PPTCourseOutlineModel.UpdateSession(ctx, session, c)
	if err != nil {
		return 0, err
	}

	return data.Id, nil
}

func (r *Render) handleLesson(ctx context.Context, session sqlx.Session, l *write_ppt.Lesson, f *types.WritePPTData) (lessonId int64, err error) {
	data, err := r.svcCtx.PPTLessonModel.FindByLessonNo(ctx, l.LessonNo)
	if err != nil {
		return 0, err
	}

	// 无题目数据
	if data == nil {
		// 传入状态为 下架，需要对其过滤
		if f.ReviewStatus == define.ReviewStatus.OffShelf {
			return 0, nil
		}
		// 传入状态为 审核通过，需要新增
		lessonId, err = r.svcCtx.PPTLessonModel.InsertSession(ctx, session, l)
		if err != nil {
			return 0, err
		}

		return lessonId, nil
	}

	// 有题目数据，强制对其进行update
	l.Id = data.Id
	l.Title = data.Title       // 保留脚本生成的数据
	l.SubTitle = data.SubTitle // 保留脚本生成的数据
	err = r.svcCtx.PPTLessonModel.UpdateSession(ctx, session, l)
	if err != nil {
		return 0, err
	}

	return data.Id, nil
}
