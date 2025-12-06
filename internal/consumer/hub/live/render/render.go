package render

import (
	"context"
	"database/sql"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"muse-admin/internal/define"
	"muse-admin/internal/model/nightstar"
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
func (r *Render) Render(ctx context.Context, f *types.LiveData) (err error) {
	var teacherData *nightstar.LiveTeacher
	// todo 不接受下架数据，需要自己干预
	if f.ReviewStatus == 5 {
		return nil
	}

	// 检查课堂编号是否存在
	dbCourse, err := r.svcCtx.LiveCourseModel.FindOneByNo(ctx, f.LiveNo)
	if err != nil {
		return err
	}
	if dbCourse == nil {
		return errors.New(fmt.Sprintf("推送的课堂:%s，课堂编号：%s，关联关系不存在，请检查", f.Name, f.LiveNo))
	}

	// 1. 先获取老师信息 ns_live_teacher 根据模式类型来决定老师信息
	if f.ModeType == 1 {
		// 模版模式
		teacherData, err = r.getTeacher(ctx, f)
	} else {
		// 定制工厂
		teacherData, err = r.getTeacherByFactory(ctx, f)
	}
	if err != nil {
		return err
	}

	// 2. 更新课堂表：ns_live_course
	c := convertCourse(f, dbCourse, teacherData.Id)
	courseId, err := r.handleCourse(ctx, c, dbCourse, f.ReviewStatus)
	if err != nil {
		return err
	}

	// 初始化transfer结构体数据
	t := &Transfer{
		CourseId:  courseId,
		TeacherId: teacherData.Id,
	}

	// 3. 大内容解析，填充课中流转所需数据
	err = r.getTransfer(ctx, f, t)
	if err != nil {
		return err
	}

	// 事务处理
	err = r.processTx(ctx, f, t)
	if err != nil {
		return err
	}

	// 将此条mq消费掉
	return nil
}

func (r *Render) getTeacher(ctx context.Context, f *types.LiveData) (*nightstar.LiveTeacher, error) {
	teacher, err := r.svcCtx.LiveTeacherModel.FindOneByTeacherCode(ctx, f.TeacherCode)
	if err != nil {
		return nil, err
	}
	if teacher == nil {
		return nil, errors.New(fmt.Sprintf("推送的课堂:%s，老师编码：%s，老师信息为空，请检查", f.Name, f.TeacherCode))
	}

	// 兜底逻辑：确保老师有讲课动图和静图
	if teacher.TeacherGif != "" || teacher.TeacherPng != "" {
		teacher.TeacherGif = f.TeacherGif
		teacher.TeacherPng = f.TeacherPng
		_ = r.svcCtx.LiveTeacherModel.Update(ctx, teacher)
	}

	return teacher, nil
}

// getTeacherByFactory 定制工厂模式下获取老师信息，若没有时就需要新增
func (r *Render) getTeacherByFactory(ctx context.Context, f *types.LiveData) (*nightstar.LiveTeacher, error) {
	t := convertTeacher(f)
	teacher, err := r.svcCtx.LiveTeacherModel.FindOneByTeacherCode(ctx, f.TeacherCode)
	if err != nil {
		return nil, err
	}

	if teacher == nil {
		// 无详情数据，做新增操作
		result, err := r.svcCtx.LiveTeacherModel.Insert(ctx, t)
		if err != nil {
			return nil, err
		}
		t.Id, err = result.LastInsertId()
		if err != nil {
			return nil, err
		}
		if t.Id == 0 {
			return nil, errors.New("新增失败，未生成自增ID")
		}
	} else {
		// 有详情数据，做修改操作
		t.Id = teacher.Id
		err = r.svcCtx.LiveTeacherModel.Update(ctx, t)
		if err != nil {
			return nil, err
		}
	}

	return t, nil
}

// 提取大内容中的课堂相关数据
func (r *Render) getTransfer(ctx context.Context, f *types.LiveData, t *Transfer) error {
	var (
		err              error
		content          Content   // 大内容解析
		frontEvent       []Message // 班主任欢迎环节
		smallTalkEvent   []Message // 寒暄环节
		smallTalkPrecast []PreCast // 寒暄开场预制
		learnEvent       []Message // 课中环节预制
		gameIds          []int64   // 当前课堂的所有事件ID
		gameStart        []float64 // 当前课堂的所有事件触发的start
	)
	err = json.Unmarshal([]byte(f.Content), &content)
	if err != nil {
		return errors.New(fmt.Sprintf("推送的课中内容，格式化失败，错误信息：%v，请检查", err))
	}
	if content.Links == nil {
		return errors.New("推送的课中内容，环节数据为空，请检查")
	}

	for _, link := range content.Links {
		var timeReset int
		// 加工时间线数据
		for k, v := range link.Timeline {
			var (
				gameId   int64
				toUserId int64
			)
			// 打开摄像头
			if v.Type == "turn_on_camera" {
				if v.CameraUserId == 0 {
					return errors.New(fmt.Sprintf("课堂编号：%s，班主任闲聊环节，需要打开摄像头的用户ID为0，请检查", f.LiveNo))
				}
				v.Type = "show-camera"
				v.Msg.UserID = v.CameraUserId
			}
			// 聊天
			if v.Type == "user_msg" {
				v.Type = "chat"
			}
			// @回复
			if len(v.Msg.Mentions) > 0 {
				toUserId = v.Msg.Mentions[0]
			}
			// 游戏环节
			if v.Type == "game_input" || v.Type == "game_select" {
				if len(v.List) <= 0 {
					return errors.New(fmt.Sprintf("课堂编号：%s，互动事件标题：%s，列表为空，请检查", f.LiveNo, v.Title))
				}
				e := convertCourseEvent(&v, t.CourseId)
				gameId, err = r.handleEvent(ctx, e)
				if err != nil {
					return err
				}
				gameIds = append(gameIds, gameId)
				gameStart = append(gameStart, v.Start) // 这里记录触发事件的时间，需要在navigate中使用，用于导航栏中所记录的个数

				// 处理题目
				// q, o := convertCourseQuestionAndOption(&v, t.CourseId, gameId)
				// r.handleCourseQuestion(ctx, gameId, q, o)
				err = r.handleCourseQuestion(ctx, t.CourseId, gameId, &v)
				if err != nil {
					return err
				}

				// 手动打上游戏标识
				v.Type = "game"
			}
			// 笔记和课后作业处理
			if v.Type == "note" || v.Type == "homework" {
				v.Type = "chat"
				v.Msg.Say = v.Text
				v.Msg.UserID = 999 // 班主任角色
			}

			msg := Message{
				Id:       cast.ToInt64(k + 1),
				Type:     v.Type,
				Time:     v.Start,
				UserId:   v.Msg.UserID,
				Content:  v.Msg.Say,
				ToUserId: toUserId, // 先取多人的第一个
				GameId:   gameId,
			}
			// 班主任闲聊环节
			if link.Type == "SMALL_CHAT" {
				t.PosterImage = link.Screen.Image.Media
				if k >= 1 {
					timeReset += util.SliceRangeRandom(f.IntervalDuration) // 需要给出一个指定范围内的随机数
				}
				msg.Time = cast.ToFloat64(timeReset)
				frontEvent = append(frontEvent, msg)
			}
			// 大寒暄环节中，需要区分寒暄环节和寒暄环节的开场预制
			if link.Type == "TALK" {
				if v.Msg.Tag == 3 {
					if k >= 1 {
						timeReset += util.SliceRangeRandom(f.IntervalDuration) // 需要给出一个指定范围内的随机数
					}
					msg.Time = cast.ToFloat64(timeReset)
					smallTalkEvent = append(smallTalkEvent, msg)
				}
				if v.Msg.Tag == 2 {
					preCast := PreCast{
						UserId:  v.Msg.UserID,
						SayType: "text",
						Say:     v.Msg.Say,
					}
					// 主讲老师的消息需要写成speak
					if v.Msg.UserID == 100 {
						preCast.SayType = "speak"
					}
					smallTalkPrecast = append(smallTalkPrecast, preCast)
				}
			}
			// 课中环节
			if link.Type == "IN_CLASS" {
				learnEvent = append(learnEvent, msg)
			}
		}

	}

	// 班主任闲聊环节
	frontEventByte, _ := json.Marshal(frontEvent)
	t.FrontEvent = string(frontEventByte)

	// 寒暄环节
	smallTalkEventByte, _ := json.Marshal(smallTalkEvent)
	t.SmallTalkEvent = string(smallTalkEventByte)

	// 寒暄环节的开场预制
	smallTalkPrecastByte, _ := json.Marshal(smallTalkPrecast)
	t.SmallTalkPrecast = string(smallTalkPrecastByte)

	// 课中环节
	learnEventByte, _ := json.Marshal(learnEvent)
	t.LearnEvent = string(learnEventByte)

	// 导航栏中的互动事件触发时间列表
	t.NavigateGameStart = gameStart

	// 清理无用事件
	r.handleCleanEvent(ctx, t.CourseId, gameIds)

	return nil
}

func (r *Render) processTx(ctx context.Context, f *types.LiveData, t *Transfer) error {
	err := r.svcCtx.MysqlConnAbility.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {

		// 1. 插入直播课堂详情表：ns_live_course_detail
		detail := convertCourseDetail(f, t)
		if err := r.handleCourseDetail(ctx, session, detail); err != nil {
			return err
		}

		// 2. 插入直播课堂时间线表：ns_live_course_timeline
		timeline := convertCourseTimeline(f, t)
		if err := r.handleCourseTimeline(ctx, session, timeline); err != nil {
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

func (r *Render) handleCourse(ctx context.Context, c *nightstar.LiveCourse, dbCourse *nightstar.LiveCourse, reviewStatus int64) (courseId int64, err error) {
	// 若该课堂已经为归档状态，则代表此条数据需要新增

	// todo 上线时，需要改为dbCourse.Status >= 3
	if dbCourse.Status > 3 {
		// 传入状态为 下架，需要对其过滤
		if reviewStatus == define.ReviewStatus.OffShelf {
			return 0, nil
		}

		// 记录版本号
		c.Version = dbCourse.Version + 1 // 版本号在上一条数据基础上+1

		// 传入状态为 审核通过，需要新增
		result, err := r.svcCtx.LiveCourseModel.Insert(ctx, c)
		if err != nil {
			return 0, err
		}
		// 获取新增ID
		courseId, err = result.LastInsertId()
		if err != nil {
			return 0, err
		}
		if courseId == 0 {
			return 0, errors.New("新增互动事件失败，未生成自增ID")
		}

		return courseId, nil
	} else {
		c.Status = 3 // todo 上线需要删除

		// 有课节数据，直接更新该课堂数据
		c.Id = dbCourse.Id
		err = r.svcCtx.LiveCourseModel.Update(ctx, c)
		if err != nil {
			return 0, err
		}
		return dbCourse.Id, nil
	}
}

func (r *Render) handleEvent(ctx context.Context, e *nightstar.LiveCourseEvent) (eventId int64, err error) {
	data, err := r.svcCtx.LiveCourseEventModel.FindOneByCourseIdAndStartId(ctx, e.LiveCourseId, e.StartId)
	if err != nil {
		return 0, err
	}

	if data == nil {
		// 新互动事件，新增
		result, err := r.svcCtx.LiveCourseEventModel.Insert(ctx, e)
		if err != nil {
			return 0, err
		}
		// 获取新增ID
		eventId, err = result.LastInsertId()
		if err != nil {
			return 0, err
		}
		if eventId == 0 {
			return 0, errors.New("新增互动事件失败，未生成自增ID")
		}
	} else {
		// 已有互动事件，修改
		e.Id, eventId = data.Id, data.Id
		err = r.svcCtx.LiveCourseEventModel.Update(ctx, e)
		if err != nil {
			return 0, err
		}
	}

	return eventId, nil
}

func (r *Render) handleCleanEvent(ctx context.Context, courseId int64, gameIds []int64) {
	// 获取课堂下所有的事件
	list, _ := r.svcCtx.LiveCourseEventModel.FindAllByCourseId(ctx, courseId)

	// 记录需要删除旧数据的startId
	var deleteEventIds []int64
	for _, v := range list {
		if !util.IsExist(v.Id, gameIds) {
			deleteEventIds = append(deleteEventIds, v.Id)
		}
	}

	// 需要清理 这次数据没有，而上次有的数据
	if len(deleteEventIds) > 0 {
		// 清理废物事件
		_ = r.svcCtx.LiveCourseEventModel.BatchDeleteByIds(ctx, deleteEventIds)
		for _, v := range deleteEventIds {
			// 清理废物事件下的题目
			question, _ := r.svcCtx.LiveCourseQuestionModel.FindAllByCourseEventId(ctx, v)
			if question == nil {
				return
			}
			var deleteQuestionIds []int64
			for _, q := range question {
				deleteQuestionIds = append(deleteQuestionIds, q.Id)
			}
			// 清理废物事件下的题目
			_ = r.svcCtx.LiveCourseQuestionModel.BatchDeleteByIds(ctx, deleteQuestionIds)
			// 清理废物事件下题目的选项
			_ = r.svcCtx.LiveCourseQuestionOptionModel.BatchDeleteByQuestionIds(ctx, deleteQuestionIds)
		}
	}
}

func (r *Render) handleCourseQuestion(ctx context.Context, courseId, eventId int64, t *Timeline) error {
	list, err := r.svcCtx.LiveCourseQuestionModel.FindAllByCourseEventId(ctx, eventId)
	if err != nil {
		return err
	}

	// 这次eventId有，上次也有，删除上次的所有数据
	var deleteIds []int64
	for _, o := range list {
		deleteIds = append(deleteIds, o.Id)
	}

	if len(deleteIds) > 0 {
		err = r.svcCtx.LiveCourseQuestionModel.BatchDeleteByIds(ctx, deleteIds)
		if err != nil {
			return err
		}
		err = r.svcCtx.LiveCourseQuestionOptionModel.BatchDeleteByQuestionIds(ctx, deleteIds)
		if err != nil {
			return err
		}
	}

	var questionType int64
	if t.Type == "game_input" {
		questionType = 1
	}
	if t.Type == "game_select" {
		questionType = 2
	}
	for k, v := range t.List {
		// 答题要求需要给一个默认兜底
		answerAsk := "回答上来一部分就算完成答题，就可以继续讲课了"
		if v.AnswerAsk != "" {
			answerAsk = v.AnswerAsk
		}
		question := &nightstar.LiveCourseQuestion{
			LiveCourseId:      courseId,
			LiveCourseEventId: eventId,
			QuestionType:      questionType,
			Question:          v.Question,
			AnswerAsk:         answerAsk,
			Tips:              v.Tips,
			Image:             v.Img,
			Answer:            v.Answer,
			Analysis: sql.NullString{
				String: v.Analyse,
				Valid:  v.Analyse != "",
			},
			Sequence: cast.ToFloat64(k + 1),
		}
		result, err := r.svcCtx.LiveCourseQuestionModel.Insert(ctx, question)
		if err != nil {
			return err
		}
		questionId, err := result.LastInsertId()
		if err != nil {
			return err
		}

		// 若题目有选项，需要关联题目ID
		if len(v.Options) > 0 {
			for optKey, optValue := range v.Options {
				option := &nightstar.LiveCourseQuestionOption{
					LiveCourseQuestionId: questionId,
					Tag:                  optValue.Tag,
					Content:              optValue.Content,
					Tips:                 optValue.Tips,
					Image:                optValue.Img,
					Sequence:             cast.ToFloat64(optKey + 1),
				}
				_, _ = r.svcCtx.LiveCourseQuestionOptionModel.Insert(ctx, option)
			}
		}
	}

	return nil
}

// func (r *Render) handleCourseQuestion(ctx context.Context, eventId int64, q []*nightstar.LiveCourseQuestion, o []*nightstar.LiveCourseQuestionOption) error {
// 	list, err := r.svcCtx.LiveCourseQuestionModel.FindAllByCourseEventId(ctx, eventId)
// 	if err != nil {
// 		return err
// 	}
//
// 	// 删除旧数据
// 	var deleteIds []string
// 	for _, o := range list {
// 		deleteIds = append(deleteIds, cast.ToString(o.Id))
// 	}
// 	err = r.svcCtx.LiveCourseQuestionModel.BatchDelete(ctx, deleteIds)
// 	if err != nil {
// 		return err
// 	}
//
// 	// 新增新数据
// 	err = r.svcCtx.LiveCourseQuestionModel.BatchInsert(ctx, q)
// 	if err != nil {
// 		return err
// 	}
//
// 	// // 旧数据
// 	// var oldMap map[int64]string
// 	// for _, o := range list {
// 	// 	oldMap[o.Id] = o.Question
// 	// }
// 	//
// 	// // 新数据对比旧数据
// 	// var deleteIds []int
// 	// for _, n := range q {
// 	// 	// 新数据等于旧数据，做更新操作
// 	// 	if n.Question ==  {
// 	//
// 	// 	} else {
// 	// 		// 将旧数据删除，新增新数据
// 	// 		deleteIds = append(deleteIds, int(o.Id))
// 	// 		err = r.svcCtx.LiveCourseQuestionModel.Insert(ctx, n)
// 	// 		if err != nil {
// 	// 			return err
// 	// 		}
// 	// 	}
// 	// }
//
// 	return nil
// }

func (r *Render) handleCourseDetail(ctx context.Context, session sqlx.Session, detail *nightstar.LiveCourseDetail) error {
	data, err := r.svcCtx.LiveCourseDetailModel.FindOneByCourseId(ctx, detail.LiveCourseId)
	if err != nil {
		return err
	}

	if data == nil {
		// 无详情数据，做新增操作
		_, err = r.svcCtx.LiveCourseDetailModel.InsertSession(ctx, session, detail)
		if err != nil {
			return err
		}
	} else {
		// 有详情数据，做修改操作
		detail.Id = data.Id
		err = r.svcCtx.LiveCourseDetailModel.UpdateSession(ctx, session, detail)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Render) handleCourseTimeline(ctx context.Context, session sqlx.Session, timeline *nightstar.LiveCourseTimeline) error {
	data, err := r.svcCtx.LiveCourseTimelineModel.FindOneByCourseId(ctx, timeline.LiveCourseId)
	if err != nil {
		return err
	}

	if data == nil {
		// 无详情数据，做新增操作
		_, err = r.svcCtx.LiveCourseTimelineModel.InsertSession(ctx, session, timeline)
		if err != nil {
			return err
		}
	} else {
		timeline.Id = data.Id
		// 有详情数据，做修改操作
		err = r.svcCtx.LiveCourseTimelineModel.UpdateSession(ctx, session, timeline)
		if err != nil {
			return err
		}
	}

	return nil
}
