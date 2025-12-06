package render

import (
	"database/sql"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"errors"
	"github.com/spf13/cast"
	"muse-admin/internal/model/kingclub/supertrain"
	"muse-admin/internal/types"
	"strconv"
	"strings"
	"time"
)

// convertCourse 转换为aw_course_outline表可录入的数据
func convertCourse(mode string, data *types.SuperTrainData) (*supertrain.Course, error) {
	r, err := util.MarshalCamelDefault(data.CourseExtra)
	if err != nil {
		return nil, err
	}

	openTime := time.Now()
	if data.OpenTime != "" {
		openTime, err = util.GetStandardDatetime(data.OpenTime)
		if err != nil {
			return nil, errors.New("开始时间日期格式错误")
		}
	}

	courseType := int64(0)
	switch data.CourseType {
	case 1:
		courseType = 2 // 超练写作
	case 2:
		courseType = 3 // 超练阅读
	}

	status := int64(2)
	// 测试环境或开发环境时，课程状态为已上架
	if mode == "test" || mode == "dev" || mode == "pro" {
		status = 3
	}

	// 默认根据序列号后三位来排序，暂定，后续可能会修改
	var sequence float64
	serial := data.CourseNo[len(data.CourseNo)-3:]
	trimmed := strings.Trim(serial, "0")           // 删除前后的 '0'
	_, errFloat := strconv.ParseFloat(trimmed, 64) // 尝试解析为 float64
	_, errInt := strconv.Atoi(trimmed)             // 尝试解析为 int
	if errFloat != nil || errInt != nil {
		sequence = cast.ToFloat64(0)
	} else {
		sequence = cast.ToFloat64(trimmed)
	}

	return &supertrain.Course{
		CourseNo:   data.CourseNo,
		UseCases:   1, // 应用场景：1.练习 2.直播互动
		CourseType: courseType,
		Version:    2, // 版本号 写死2
		Name:       data.CourseName,
		LessonName: data.LessonName,
		Subject:    data.Subject,
		Unit:       data.Unit,
		Image:      data.Image,
		Intro:      data.Intro,
		OpenTime:   openTime,
		Type:       data.Type,
		Status:     status,
		Level: sql.NullString{
			String: "{\"11\": [\"1级\", \"1级\", \"1级\", \"1级\", \"1级\", \"1级\", \"1级\"], \"12\": [\"2级\", \"2级\", \"2级\", \"2级\", \"2级\", \"2级\", \"2级\"], \"13\": [\"3级\", \"3级\", \"3级\", \"3级\", \"3级\", \"3级\", \"3级\"], \"14\": [\"4级\", \"4级\", \"4级\", \"4级\", \"4级\", \"4级\", \"4级\"], \"15\": [\"5级\", \"5级\", \"5级\", \"5级\", \"5级\", \"5级\", \"5级\"], \"16\": [\"6级\", \"6级\", \"6级\", \"6级\", \"6级\", \"6级\", \"6级\"], \"21\": [\"7级\", \"7级\", \"7级\", \"7级\", \"7级\", \"7级\", \"7级\"], \"22\": [\"8级\", \"8级\", \"8级\", \"8级\", \"8级\", \"8级\", \"8级\"], \"23\": [\"9级\", \"9级\", \"9级\", \"9级\", \"9级\", \"9级\", \"9级\"], \"31\": [\"10级\", \"10级\", \"10级\", \"10级\", \"10级\", \"10级\", \"10级\"], \"32\": [\"11级\", \"11级\", \"11级\", \"11级\", \"11级\", \"11级\", \"11级\"], \"33\": [\"12级\", \"12级\", \"12级\", \"12级\", \"12级\", \"12级\", \"12级\"]}",
			Valid:  true,
		},
		Remark: "",
		Extra: sql.NullString{
			String: string(r),
			Valid:  true,
		},
		Sequence: sequence,
	}, nil
}

// convertChapter 转换为ac_course_chapter表可录入的数据
func convertChapter(courseType int64, data *types.Chapter, courseId int64) (*supertrain.CourseChapter, error) {
	// 如果是阅读课程，并且引导视频不为空，则将引导视频赋值给引导视频地址
	if courseType == 2 && data.Extra.GuideVideo != "" {
		data.Extra.GuideVideoUrl = data.Extra.GuideVideo
	}
	var (
		r   []byte
		err error
	)
	// 写作课程，需要把extra中所有的字段作为小驼峰存储
	if courseType == 1 {
		r, err = util.MarshalCamelAllForceCamel(data.Extra)
	} else {
		r, err = util.MarshalCamelDefault(data.Extra)
	}
	if err != nil {
		return nil, err
	}
	// isNew, isLock := int64(1), int64(2)
	// if data.IsNew {
	// 	isNew = 2
	// }
	// if data.IsLock {
	// 	isLock = 1
	// }
	return &supertrain.CourseChapter{
		ChapterNo:  data.No,
		Type:       data.Type,
		CourseId:   courseId,
		Name:       data.Name,
		Title:      data.Title,
		Image:      data.Image,
		Intro:      data.Intro,
		Index:      1, // 经听森沟通后，暂时写死为1
		GuideVideo: data.GuideVideo,
		Sequence:   data.Sequence,
		TeacherId:  data.TeacherId,
		Status:     3,
		CanLearn:   1, // 是否能学 1.可以 2.不可以
		IsNew:      2, // 是否上新 1.否 2.是
		Extra: sql.NullString{
			String: string(r),
			Valid:  true,
		},
	}, nil
}

// convertTask 转换为ac_chapter_task表可录入的数据
func convertTask(courseType int64, data *types.Task, courseId, chapterId int64) (*supertrain.ChapterTask, error) {
	// 默认按照段落类型来处理
	var paragraphType int64
	if data.Extra.ParagraphType {
		paragraphType = 1 // 按字数计分
	} else {
		paragraphType = 2 // 按任务计分
	}
	taskExtra := TaskExtra{
		ParagraphType: paragraphType,
	}

	var (
		r   []byte
		err error
	)
	// 写作课程，需要把extra中所有的字段作为小驼峰存储
	if courseType == 1 {
		r, err = util.MarshalCamelAllForceCamel(taskExtra)
	} else {
		r, err = util.MarshalCamelDefault(taskExtra)
	}
	if err != nil {
		return nil, err
	}

	return &supertrain.ChapterTask{
		TaskNo:          data.TaskNo,
		Title:           data.Title,
		Subtitle:        data.Subtitle,
		Description:     data.Description,
		CourseId:        courseId,
		CourseChapterId: chapterId,
		ProgressPrefix:  data.ProgressPrefix,
		Extra:           string(r),
		Type:            data.Type,
		Sequence:        data.Sequence,
		Deleted:         1, // 1.未删除 2.已删除
	}, nil
}

func convertSubTask(courseType int64, paragraphType bool, data *types.SubTask, taskId int64) (*supertrain.SubTask, error) {
	// 若根据字数来计分，需要确保 GradeRule 不能为空
	if paragraphType && len(data.Extra.GradeRules) <= 0 {
		return nil, errors.New("grade rule is empty")
	}

	var (
		r   []byte
		err error
	)
	// 写作课程，需要把extra中所有的字段作为小驼峰存储
	if courseType == 1 {
		r, err = util.MarshalCamelAllForceCamel(data.Extra)
	} else {
		r, err = util.MarshalCamelDefault(data.Extra)
	}
	if err != nil {
		return nil, err
	}

	return &supertrain.SubTask{
		Mode:          data.Mode,
		SubTaskNo:     data.SubTaskNo,
		ChapterTaskId: taskId,
		Title:         data.Title,
		Article: sql.NullString{
			String: data.Article,
			Valid:  true,
		},
		Description: data.Description,
		ImageUrl:    data.ImageUrl,
		Keywords:    "[]", // 废弃字段
		Extra:       string(r),
		Type:        data.Type,
		Sequence:    data.Sequence,
		Deleted:     1,
	}, nil
}

func convertMetal(data types.Medal) (*supertrain.Medal, error) {
	var description string
	if data.Description == "" {
		description = "阅读任务"
	}
	return &supertrain.Medal{
		Name:            data.Name,
		Description:     description,
		Type:            data.Type,
		PromptKnowledge: data.PromptKnowledge,
	}, nil
}
