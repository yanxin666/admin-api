package render

import (
	"database/sql"
	"encoding/json"
	"github.com/spf13/cast"
	"muse-admin/internal/model/nightstar"
	"muse-admin/internal/types"
	"sort"
)

// convertTeacher 转换为ns_live_teacher表可录入的数据
func convertTeacher(f *types.LiveData) *nightstar.LiveTeacher {
	return &nightstar.LiveTeacher{
		TeacherCode:      f.TeacherCode,
		TeacherName:      f.TeacherName,
		TeacherNickname:  f.TeacherNickname,
		HandAudio:        f.HandAudio,
		ToneType:         f.TeacherToneType,
		TeacherTone:      f.TeacherTone,
		ToneModel:        f.TeacherToneModel,
		TeacherIntroduce: f.TeacherIntroduce,
		TeacherGif:       f.TeacherGif,
		TeacherPng:       f.TeacherPng,
	}
}

// convertCourse 转换为ns_live_course表可录入的数据
func convertCourse(f *types.LiveData, dbCourse *nightstar.LiveCourse, teacherId int64) *nightstar.LiveCourse {
	course := &nightstar.LiveCourse{
		Id:        dbCourse.Id,
		LiveNo:    f.LiveNo,
		Version:   dbCourse.Version, // 默认取上次的版本号
		Name:      dbCourse.Name,
		TeacherId: teacherId,
		Label:     dbCourse.Label,
		Image:     dbCourse.Image,
		Intro:     dbCourse.Intro,
		Keyword:   dbCourse.Keyword,
		// OpenTime:  time.Unix(0, 0), // 默认1970
		OpenTime: dbCourse.OpenTime,
		Type:     dbCourse.Type, // 默认取上次的课堂类型
		Status:   2,             // 默认审核中
		Remark:   dbCourse.Remark,
		Sequence: dbCourse.Sequence,
	}

	// 预览时可以干预名称
	if f.EnvType == 1 {
		course.Name = f.Name
	}

	return course
}

// convertCourseDetail 转换为ns_live_course_detail表可录入的数据
func convertCourseDetail(f *types.LiveData, t *Transfer) *nightstar.LiveCourseDetail {
	return &nightstar.LiveCourseDetail{
		LiveCourseId: t.CourseId,
		TeacherId:    t.TeacherId,
		TeacherGif:   f.TeacherGif, // 每堂课老师动图与静图可能不一样，顾记录每次课堂老师的讲课图片
		TeacherPng:   f.TeacherPng,
		TeacherVideo: f.TeacherVideo,
		Duration:     f.Duration,
		Topic:        f.Topic,
		PptVideo:     f.PPTVideo,
		FullImage:    f.FullImage,
		PosterImage:  t.PosterImage,
		Note:         f.Note,
		Homework:     f.Homework,
	}
}

// convertCourseTimeline 转换为ns_live_course_timeline表可录入的数据
func convertCourseTimeline(f *types.LiveData, t *Transfer) *nightstar.LiveCourseTimeline {
	// 替换字幕中的key由下划线转为驼峰
	var o []ContentDraft
	_ = json.Unmarshal([]byte(f.Draft), &o)
	draft := make([]Draft, len(o))
	for i, item := range o {
		draft[i] = Draft{
			Text:      item.Text,
			TimeBegin: item.TimeBegin,
			TimeEnd:   item.TimeEnd,
		}
	}
	draftByte, _ := json.Marshal(draft)

	// 获取每个导航栏的时间区间里的事件个数
	gameCount := countInRangesOptimized(f.Navigate, t.NavigateGameStart)

	// 替换导航栏中的key由下划线转为驼峰 && 新增该环节中互动事件的个数
	nav := make([]Navigate, len(f.Navigate))
	for i, item := range f.Navigate {
		nav[i] = Navigate{
			Name:      item.Name,
			StartTime: item.StartTime,
			Num:       gameCount[i],
		}
	}
	navByte, _ := json.Marshal(nav)

	return &nightstar.LiveCourseTimeline{
		LiveCourseId:     t.CourseId,
		Navigate:         string(navByte),
		Draft:            string(draftByte),
		DraftRag:         f.DraftRag,
		SimpleDraft:      f.SimpleDraft,
		BlockTime:        f.BlockTime,
		FrontEvent:       t.FrontEvent,
		SmallTalkEvent:   t.SmallTalkEvent,
		SmallTalkPrecast: t.SmallTalkPrecast,
		LearnEvent:       t.LearnEvent,
	}
}

// convertCourseEvent 转换为ns_live_course_event表可录入的数据
func convertCourseEvent(t *Timeline, liveCourseId int64) *nightstar.LiveCourseEvent {
	var (
		eventType  int64
		appearance int64
		promptCode string
	)
	if t.Type == "game_input" {
		// 题型风格：1.我说上句你来接 2.排序题 3.课前简答题 4.课后简答题 5.文本单选题 6.图片选择题 7.火眼金睛 8.喜欢的字 9.你的称号我知道
		eventType = 1
		switch t.Appearance {
		// 我说上句你来接
		case 1:
			appearance = 1
			promptCode = "dui_shi_flow"
		// 排序题
		case 2:
			appearance = 2
			promptCode = "sorting_question_flow"
		// 课前简答题
		case 3:
			appearance = 3
			promptCode = "short_answer_before_flow"
		// 课后简答题
		case 4:
			appearance = 4
			promptCode = "short_answer_after_flow"
		// 喜欢的字
		case 5:
			appearance = 8
			promptCode = "favorite_word_flow"
		// 你的称号我知道
		case 6:
			appearance = 9
			promptCode = "class_interaction_flow"
		// 英语：开放式问答
		case 7:
			appearance = 10
			promptCode = "dui_shi_en_flow"
		// 英语：单词跟读
		case 8:
			appearance = 11
			promptCode = "word.pronounce.en.flow"
		// 英语：练习1
		case 9:
			appearance = 12
			promptCode = "practice1.en.flow"
		// 英语：练习2
		case 10:
			appearance = 13
			promptCode = "practice2.en.flow"
		}
	}
	if t.Type == "game_select" {
		eventType = 2
		switch t.Appearance {
		// 文本单选题
		case 1:
			appearance = 5
			promptCode = "common_single_choose_flow"
		// 图片选择题
		case 2:
			appearance = 6
			promptCode = "common_single_choose_flow"
		// 火眼金睛
		case 3:
			appearance = 7
			promptCode = "fire_eyes_question_flow"
		}
	}

	var (
		talksStr = "[]"     // 互动默认值
		talks    []TaskTalk // 互动列表
	)
	if len(t.Talks) > 0 {
		// 事件列表
		for _, v := range t.Talks {
			talks = append(talks, TaskTalk{
				Say:    v.Say,
				UserID: v.UserID,
			})
		}
		r, _ := json.Marshal(talks)
		talksStr = string(r)
	}

	return &nightstar.LiveCourseEvent{
		LiveCourseId:       liveCourseId,
		StartId:            t.Start,
		EventType:          eventType,
		Appearance:         appearance,
		PromptCode:         promptCode,
		Title:              t.Title,
		Task:               t.Task,
		Talk:               talksStr,
		WithoutAnswerImage: t.WithoutAnswerImage,
		WithinAnswerImage:  t.WithinAnswerImage,
	}
}

// convertCourseQuestionAndOption 转换为ns_live_course_question表 和 ns_live_course_question_option表可录入的数据
func convertCourseQuestionAndOption(t *Timeline, liveCourseId, eventId int64) ([]*nightstar.LiveCourseQuestion, []*nightstar.LiveCourseQuestionOption) {
	var (
		questionType int64
		questionList []*nightstar.LiveCourseQuestion
		optionList   []*nightstar.LiveCourseQuestionOption
	)
	if t.Type == "game_input" {
		questionType = 1
	}
	if t.Type == "game_select" {
		questionType = 2
	}
	for k, v := range t.List {
		if len(v.Options) > 0 {
			for optKey, optValue := range v.Options {
				optionList = append(optionList, &nightstar.LiveCourseQuestionOption{
					Tag:      optValue.Tag,
					Content:  optValue.Content,
					Tips:     optValue.Tips,
					Image:    optValue.Img,
					Sequence: cast.ToFloat64(optKey + 1),
				})
			}
		}
		questionList = append(questionList, &nightstar.LiveCourseQuestion{
			LiveCourseId:      liveCourseId,
			LiveCourseEventId: eventId,
			QuestionType:      questionType,
			Question:          v.Question,
			Tips:              v.Tips,
			Image:             v.Img,
			Answer:            v.Answer,
			Analysis: sql.NullString{
				String: v.Analyse,
				Valid:  v.Analyse != "",
			},
			Sequence: cast.ToFloat64(k + 1),
		})
	}

	return questionList, optionList
}

func countInRangesOptimized(ranges []types.Navigate, arr []float64) []int64 {
	// 初始化结果数组，长度与 ranges 相同
	counts := make([]int64, len(ranges))

	// ranges[i].StartTime为秒，但num为毫秒，顾此逻辑中的StartTime需要乘1000
	for _, num := range arr {
		// 使用二分查找确定 num 属于哪个区间
		idx := sort.Search(len(ranges), func(i int) bool {
			return ranges[i].StartTime*1000 >= num
		}) - 1

		// 如果 idx >= 0，说明 ranges[idx].StartTime < num ≤ ranges[idx+1].StartTime
		// 最后一个区间是 (ranges[-1].StartTime, +∞)
		if idx >= 0 {
			if idx < len(ranges)-1 {
				if num <= ranges[idx+1].StartTime*1000 { // 左开右闭：(StartTime, EndTime]
					counts[idx]++
				}
			} else { // 最后一个区间：(StartTime, +∞)
				counts[idx]++
			}
		}
		// 如果 idx < 0，说明 num <= ranges[0].StartTime，不统计
	}

	return counts
}
