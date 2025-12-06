package render

import (
	"encoding/json"
	"muse-admin/internal/model/write_ppt"
	"muse-admin/internal/types"
	"strings"
)

// convertCourse 转换为aw_course_outline表可录入的数据
func convertCourse(data *types.WritePPTData) *write_ppt.CourseOutline {
	return &write_ppt.CourseOutline{
		Name:     "", // 脚本生成
		Unit:     data.Unit,
		Category: data.LessonCategory,
		Topic:    strings.TrimSpace(data.Topic),
		Series:   data.Series,
		Status:   1,
	}
}

// convertLesson 转换为aw_lesson表可录入的数据
func convertLesson(data *types.WritePPTData, courseId int64) *write_ppt.Lesson {
	pptUrls, _ := json.Marshal(data.PPTFiles)
	return &write_ppt.Lesson{
		CourseOutlineId:   courseId,
		LessonNo:          data.Id,
		Title:             data.Title,
		SubTitle:          "", // 脚本生成
		LessonNumber:      data.LessonNumber,
		LessonType:        data.LessonType,
		LessonContent:     strings.TrimSpace(data.LessonContent),
		LessonContentText: strings.TrimSpace(data.LessonContentText),
		Notes:             strings.TrimSpace(data.Notes),
		PptFiles:          string(pptUrls),
		MindSlices:        data.MindSlices,
		MindFull:          strings.TrimSpace(data.MindFull),
		AudioUrl:          strings.TrimSpace(data.AudioUrl),
		ReviewStatus:      data.ReviewStatus,
	}
}
