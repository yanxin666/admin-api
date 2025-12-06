package render

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/third/tencent"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"muse-admin/internal/define"
	knowledgeModel "muse-admin/internal/model/knowledge"
	"muse-admin/internal/model/lesson"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
)

type Schedule struct {
	svcCtx *svc.ServiceContext

	parentLessonData *knowledgeModel.Lesson
	cos              *tencent.Cos
}

func NewSchedule(svcCtx *svc.ServiceContext) *Schedule {

	return &Schedule{
		svcCtx: svcCtx,
	}
}

// Render 处理数据
func (s *Schedule) Render(ctx context.Context, f *types.ScheduleData) (err error) {
	s.cos = tencent.NewCos(ctx, tencent.CocConf{
		SecretId:  s.svcCtx.Config.Oss.SecretId,
		SecretKey: s.svcCtx.Config.Oss.SecretKey,
		Appid:     s.svcCtx.Config.Oss.Appid,
		Bucket:    s.svcCtx.Config.Oss.Bucket,
		Region:    s.svcCtx.Config.Oss.Region,
	})

	// 获取课程大纲列表
	courseDataArr, err := s.getCourseDataArr(ctx, f)
	if err != nil {
		return err
	}

	// 执行每个大纲下所对应的课程
	for _, courseData := range courseDataArr {
		err = s.processSingleCourse(ctx, f, courseData)
		if err != nil {
			return err
		}
	}

	// 将此条mq消费掉
	return nil
}

func (s *Schedule) getCourseDataArr(ctx context.Context, f *types.ScheduleData) ([]lesson.CourseOutline, error) {
	var (
		courseDataArr []lesson.CourseOutline
		err           error
	)

	// ID为488 弟子规, 493从家庭到国家的修身之道 这两个大语文没有父类关系，所以需要额外处理
	if f.NodeType == define.NodeType.BigLanguage && util.IsExist(f.Id, []int64{488, 493}) {
		// 验证是否具有关联关系
		courseDataArr, err = s.svcCtx.CourseOutlineModel.FindAllByMarkId(ctx, f.Id)
		if err != nil {
			return nil, err
		}
		if courseDataArr == nil {
			return nil, errors.New(fmt.Sprintf("%s，找不到关联关系，当前名称：%s，当前ID：%d", define.LessonTypeMap[f.NodeType], f.CourseName, f.Id))
		}
	}

	// 大语文 按parentId来做关联关系
	if f.NodeType == define.NodeType.BigLanguage && !util.IsExist(f.Id, []int64{488, 493}) {
		if f.ParentId == 0 {
			return nil, errors.New("大语文推送时，父类ID为0，请检查")
		}
		// 查询父类数据是否存在
		s.parentLessonData, err = s.svcCtx.LessonModel.FindOneByNo(ctx, cast.ToString(f.ParentId))
		if err != nil {
			return nil, err
		}
		// 父类数据不存在就抛错报警
		if s.parentLessonData == nil {
			// 父类数据不存在 && 数据为下架时，直接过滤
			if f.Status == define.ReviewStatus.OffShelf {
				return nil, nil
			}
			return nil, errors.New(fmt.Sprintf("%s，父类数据为空，当前传入的Id:%d,ParentId：%d", define.LessonTypeMap[f.NodeType], f.Id, f.ParentId))
		}
		// 大语文没有编号
		courseDataArr = append(courseDataArr, lesson.CourseOutline{})
	}

	// 小语文 按 f.Id 取关联关系
	if f.NodeType == define.NodeType.SmallLanguage {
		if f.Grade == 0 {
			return nil, errors.New("小语文推送时，年级为0，请检查")
		}
		// 验证是否具有关联关系
		courseDataArr, err = s.svcCtx.CourseOutlineModel.FindAllByMarkId(ctx, f.Id)
		if err != nil {
			return nil, err
		}

		// 当markId=0时，可能为旧版本的数据，需要兼容
		if courseDataArr == nil {
			// 若按remark取不到内容，就代表无关联关系
			courseDataArr, err = s.svcCtx.CourseOutlineModel.FindAllByRemark(ctx, f.CourseName)
			if courseDataArr == nil {
				return nil, errors.New(fmt.Sprintf("%s，找不到关联关系，当前名称：%s，当前ID：%d", define.LessonTypeMap[f.NodeType], f.CourseName, f.Id))
			}
		}
	}

	// 小灶课 按 f.Id 取关联关系，且目前只有等级为2的数据
	if f.NodeType == define.NodeType.SpecialCourse {
		// 验证是否具有关联关系
		courseDataArr, err = s.svcCtx.CourseOutlineModel.FindAllByMarkId(ctx, f.Id)
		if err != nil {
			return nil, err
		}

		// 当markId=0时，可能为旧版本的数据，需要兼容
		if courseDataArr == nil {
			courseDataArr, err = s.svcCtx.CourseOutlineModel.FindAllByRemark(ctx, f.PointName)
			// 若按remark取不到内容，就代表无关联关系
			if courseDataArr == nil {
				return nil, errors.New(fmt.Sprintf("%s，找不到关联关系，当前名称：%s，当前ID：%d", define.LessonTypeMap[f.NodeType], f.PointName, f.Id))
			}
		}

		if f.Level != 2 {
			return nil, errors.New(fmt.Sprintf("%s，暂不支持当前等级：%d，当前名称：%s，当前ID：%d", define.LessonTypeMap[f.NodeType], f.Level, f.PointName, f.Id))
		}
	}

	return courseDataArr, nil
}

func (s *Schedule) processSingleCourse(ctx context.Context, f *types.ScheduleData, courseData lesson.CourseOutline) error {
	// 小语文保证只出一个
	if f.NodeType == define.NodeType.SmallLanguage && f.Grade != courseData.Grade {
		return nil
	}

	// 开启事务：
	// 1. 插入课节表：kn_lesson 获取课节ID
	// 2. 插入课节知识点表：kn_lesson_point
	// 3. 插入例题表：kn_example 获取例题ID
	// 4. 插入素材表：kn_material 获取素材ID
	// 5. 插入题库表：kn_lesson_question 获取题目ID
	// 6. 插入题库选项表：kn_lesson_question_option
	// 7. 插入课节资源表：kn_lesson_resource
	err := s.svcCtx.MysqlConnAbility.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		var (
			err         error
			lessonId    int64
			exampleId   int64
			questionArr []string
		)

		// 1. 插入课节表：kn_lesson 获取课节ID
		l := convertLesson(f, courseData.LessonGroupNo) // 转换为可直接入库的数据
		lessonId, err = s.handleLesson(ctx, session, courseData, l, f)
		if err != nil {
			return err
		}
		if lessonId == 0 {
			// 数据过滤
			if err == nil {
				return nil
			} else {
				return errors.New("lessonId为空")
			}
		}

		// 2. 插入课节知识点表：kn_lesson_point
		point := convertLessonPoint(f, lessonId)
		if err = s.handleLessonPoint(ctx, session, point); err != nil {
			return err
		}

		// 3. 插入例题表：kn_example 获取例题ID
		e := convertExample(f, courseData.Title, courseData.SubTitle) // 转换为可直接入库的数据
		exampleId, err = s.handleExample(ctx, session, e)
		if err != nil {
			return err
		}

		// 4. 插入素材表：kn_material 获取素材ID
		// 5. 插入题库表：kn_lesson_question 获取题目ID
		// 6. 插入题库选项表：kn_lesson_question_option
		// 7. 插入题库tts表：kn_lesson_question_tts
		// 8. 插入课节资源表：kn_lesson_resource
		for _, v := range f.Question {
			var (
				materialId int64
				questionId int64
				isDelete   bool
			)

			// 已经执行过的不重复执行
			if util.IsExist(v.QuestionNo, questionArr) {
				continue
			}

			// 4. 插入素材表：kn_material 获取素材ID
			m := convertMaterial(&v) // 转换为可直接入库的数据
			materialId, err = s.handleMaterial(ctx, session, v, m)
			if err != nil {
				return err
			}

			// 5. 插入题库表：kn_lesson_question 获取题目ID
			// 查询题目是否存在
			q := convertQuestion(&v, f.NodeType, materialId, exampleId)
			questionId, isDelete, err = s.handleQuestion(ctx, session, v, q)
			if err != nil {
				return err
			}
			// 下架的数据需要过滤
			if questionId == 0 && err == nil {
				continue
			}

			// 6. 有题目选项
			opt := convertOption(&v, questionId)
			if v.Type == SingleChoice && opt == nil {
				return errors.New(fmt.Sprintf("当前知识点名称：%s，当前ID：%d，题目编号：%s，单选题未有选项", f.PointName, f.Id, v.QuestionNo))
			}
			if err = s.handleQuestionOpt(ctx, session, opt); err != nil {
				return err
			}

			// 7. 插入tts表
			if v.Ask != "" {
				tts, err := convertLessonQuestionTTS(ctx, s.cos, questionId, v.Ask, s.svcCtx.Config.MiniMax, s.svcCtx.Config.Oss.ReplaceDomain)
				if err = s.handleQuestionTTS(ctx, session, tts); err != nil {
					return err
				}
			}

			// 8. 操作课节资源表：kn_lesson_resource
			if err = s.handleResource(ctx, session, f, lessonId, questionId, isDelete); err != nil {
				return err
			}

			// 记录处理的题库No
			questionArr = append(questionArr, v.QuestionNo)
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

func (s *Schedule) handleLesson(ctx context.Context, session sqlx.Session, courseData lesson.CourseOutline, l *knowledgeModel.Lesson, f *types.ScheduleData) (lessonId int64, err error) {
	var lessonData *knowledgeModel.Lesson
	// 非小灶课处理逻辑
	if f.NodeType != define.NodeType.SpecialCourse {
		// 查询课节编号 lessonNo 是否存在
		lessonData, err = s.svcCtx.LessonModel.FindOneByNo(ctx, l.LessonNo)
		if err != nil {
			return 0, err
		}
		// 无课节数据
		if lessonData == nil {
			// 传入状态为 下架，需要对其过滤
			if f.Status == define.ReviewStatus.OffShelf {
				return 0, nil
			}
			// 传入状态为 审核通过，需要新增
			lessonId, err = s.svcCtx.LessonModel.InsertSession(ctx, session, l)
			if err != nil {
				return 0, err
			}
		}
		// 有课节数据，强制对其进行update
		if lessonData != nil {
			lessonId = lessonData.Id
			l.Id = lessonData.Id
			err = s.svcCtx.LessonModel.UpdateSession(ctx, session, l)
			if err != nil {
				return 0, err
			}
		}

		return lessonId, nil
	}

	// 小灶课单独处理
	// 查询 lessonGroupNo 是否存在
	s.parentLessonData, err = s.svcCtx.LessonModel.FindOneByParentLesson(ctx, courseData.LessonGroupNo, f.Level)
	if err != nil {
		return 0, err
	}
	// 首次推送过来的小灶课
	if s.parentLessonData == nil {
		// 传入状态为 下架，需要对其过滤
		if f.Status == define.ReviewStatus.OffShelf {
			return 0, nil
		}
		// 传入状态为 审核通过，需要新增
		l.Remark = f.PointName
		l.Name = f.CourseName
		lessonId, err = s.svcCtx.LessonModel.InsertSession(ctx, session, l)
		if err != nil {
			return 0, err
		}
	}
	// 非首次的小灶课，代表相同groupNo的课节数据，需要对其进行关联关系的绑定
	if s.parentLessonData != nil {
		// // 当前数据为下架，需要对首次的课节做下架操作
		// if f.Status == define.ReviewStatus.OffShelf {
		//	parentLessonData.ReviewStatus = define.ReviewStatus.OffShelf
		//	err = s.svcCtx.LessonModel.UpdateFieldsWithTx(ctx, session, parentLessonData.LessonGroupNo, parentLessonData)
		//	if err != nil {
		//		return 0,err
		//	}
		// }
		// 查询当前group下面的lessonNo课程是否存在
		lessonData, err = s.svcCtx.LessonModel.FindOneByNoAndGroupNo(ctx, l.LessonNo, l.LessonGroupNo)
		if err != nil {
			return 0, err
		}
		// 无课节数据
		if lessonData == nil {
			// 传入状态为 下架，需要对其过滤
			if f.Status == define.ReviewStatus.OffShelf {
				return 0, nil
			}
			// 对非首次的小灶课做一个数据关联的记录
			l.ParentId = cast.ToInt64(s.parentLessonData.LessonNo)
			l.Remark = courseData.Remark
			l.Name = f.CourseName
			lessonId, err = s.svcCtx.LessonModel.InsertSession(ctx, session, l)
			if err != nil {
				return 0, err
			}
		}
		// 有课节数据，强制对其进行update
		if lessonData != nil {
			lessonId = lessonData.Id
			l.Id = lessonData.Id
			// 改数据之前保持之前特殊处理的数据
			l.ParentId = lessonData.ParentId
			l.Remark = lessonData.Remark
			l.Name = lessonData.Name
			err = s.svcCtx.LessonModel.UpdateSession(ctx, session, l)
			if err != nil {
				return 0, err
			}
		}
	}

	return lessonId, nil
}

func (s *Schedule) handleLessonPoint(ctx context.Context, session sqlx.Session, point []knowledgeModel.LessonPoint) error {
	for _, op := range point {
		// 查询课节资源是否存在
		pointId, err := s.svcCtx.LessonPointModel.FindResourceId(ctx, &op)
		if err != nil {
			return err
		}
		// 若查询的资源ID不存在，新增资源表
		if pointId == 0 {
			_, err = s.svcCtx.LessonPointModel.InsertSession(ctx, session, &op)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Schedule) handleExample(ctx context.Context, session sqlx.Session, e *knowledgeModel.Example) (exampleId int64, err error) {
	if e.Explain != "" || e.PlainText != "" || e.Title != "" || e.SubTitle != "" || e.LessonNotes != "" {
		// 查询例题ID是否存在
		exampleId, err = s.svcCtx.ExampleModel.FindExampleId(ctx, e)
		if err != nil {
			return 0, err
		}
		// 若查询的例题ID不存在，新增例题表
		if exampleId == 0 {
			exampleId, err = s.svcCtx.ExampleModel.InsertSession(ctx, session, e)
			if err != nil {
				return 0, err
			}
		}
	}

	return exampleId, nil
}

func (s *Schedule) handleMaterial(ctx context.Context, session sqlx.Session, v types.ScheduleDataQuestion, m *knowledgeModel.Material) (materialId int64, err error) {
	if m.Title != "" || m.Author != "" || m.Source != "" || m.Content.String != "" || m.Background.String != "" || m.AuthorIntro.String != "" {
		// 查询素材ID是否存在
		materialId, err = s.svcCtx.MaterialModel.FindMaterialId(ctx, m)
		if err != nil {
			return 0, err
		}
		// 若查询的素材ID不存在，新增素材表
		if materialId == 0 {
			// 若当前题状态为 下架，则不用新增素材表
			if v.ReviewStatus == define.ReviewStatus.OffShelf {
				return 0, nil
			}
			materialId, err = s.svcCtx.MaterialModel.InsertSession(ctx, session, m)
			if err != nil {
				return 0, err
			}
		}
	}

	return materialId, nil
}

func (s *Schedule) handleQuestion(ctx context.Context, session sqlx.Session, v types.ScheduleDataQuestion, q *knowledgeModel.LessonQuestion) (questionId int64, isDelete bool, err error) {
	questionData, err := s.svcCtx.LessonQuestionModel.FindByQuestionNo(ctx, v.QuestionNo)
	if err != nil {
		return 0, isDelete, err
	}

	// 无题目数据
	if questionData == nil {
		// 传入状态为 下架，需要对其过滤
		if v.ReviewStatus == define.ReviewStatus.OffShelf {
			return 0, isDelete, nil
		}
		// 传入状态为 审核通过，需要新增
		questionId, err = s.svcCtx.LessonQuestionModel.InsertSession(ctx, session, q)
		if err != nil {
			return 0, isDelete, err
		}
	}
	// 有题目数据，强制对其进行update
	if questionData != nil {
		// 传入状态 == 数据库状态，代表其数据未做改动
		// if v.ReviewStatus != questionData.ReviewStatus {
		//	continue
		// }
		// 传入状态 != 数据库状态，代表其状态有改动，需要对其进行update
		// if v.ReviewStatus != questionData.ReviewStatus {
		questionId = questionData.Id
		q.Id = questionData.Id
		err = s.svcCtx.LessonQuestionModel.UpdateSession(ctx, session, q)
		if err != nil {
			return 0, isDelete, err
		}
		// }

		if v.ReviewStatus == define.ReviewStatus.OffShelf {
			isDelete = true
		}
	}

	return questionId, isDelete, nil
}

func (s *Schedule) handleQuestionOpt(ctx context.Context, session sqlx.Session, opt []knowledgeModel.LessonQuestionOption) error {
	for _, op := range opt {
		_, err := s.svcCtx.LessonQuestionOptionModel.InsertSession(ctx, session, &op)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Schedule) handleQuestionTTS(ctx context.Context, session sqlx.Session, tts *knowledgeModel.LessonQuestionTts) error {
	data, err := s.svcCtx.LessonQuestionTtsModel.FindByQuestionId(ctx, tts.QuestionId)
	if err != nil {
		return err
	}

	// 无题目数据
	if data == nil {
		_, err = s.svcCtx.LessonQuestionTtsModel.InsertSession(ctx, session, tts)
		if err != nil {
			return err
		}
	}

	// 有题目数据，进行更新
	if data != nil {
		tts.Id = data.Id
		err = s.svcCtx.LessonQuestionTtsModel.UpdateSession(ctx, session, tts)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Schedule) handleResource(ctx context.Context, session sqlx.Session, f *types.ScheduleData, lessonId, questionId int64, isDelete bool) error {
	var selfLessonId int64 // 本身的子类ID，方便删除本身子类数据时使用
	selfLessonId = lessonId
	// 大语文 和 小灶课 额外处理，若有父类，需要关联到父类里
	if s.parentLessonData != nil {
		lessonId = s.parentLessonData.Id
	}
	// 查询课节资源是否存在
	resourceId, err := s.svcCtx.LessonResourceModel.FindOneByLessonIdAndQuestionId(ctx, lessonId, questionId)
	if err != nil {
		return err
	}
	// 若此条数据为下线状态，需要删除所下线题目的关联关系
	if isDelete {
		if resourceId == 0 {
			return nil
		}
		// 删除父类的关联关系
		_ = s.svcCtx.LessonResourceModel.DeleteSession(ctx, session, resourceId)
		// 删除子类的关联关系
		selfResourceId, _ := s.svcCtx.LessonResourceModel.FindOneByLessonIdAndQuestionId(ctx, selfLessonId, questionId)
		if selfResourceId != 0 {
			_ = s.svcCtx.LessonResourceModel.DeleteSession(ctx, session, selfResourceId)
		}
		return nil
	}

	// 非下线数据，存在则返回
	if resourceId != 0 {
		return nil
	}

	// 若查询的资源ID不存在，新增资源表
	if lessonId == selfLessonId {
		// 新增数据
		resourceId, err = s.svcCtx.LessonResourceModel.InsertSession(ctx, session, convertLessonResource(lessonId, questionId, f))
	} else {
		// 父类新增
		resourceId, err = s.svcCtx.LessonResourceModel.InsertSession(ctx, session, convertLessonResource(lessonId, questionId, f))
		// 子类新增
		resourceId, err = s.svcCtx.LessonResourceModel.InsertSession(ctx, session, convertLessonResource(selfLessonId, questionId, f))
	}
	return nil
}
