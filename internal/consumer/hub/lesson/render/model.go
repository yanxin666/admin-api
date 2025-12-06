package render

import (
	"context"
	"database/sql"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/third/minimax"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/third/tencent"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logc"
	"muse-admin/internal/config"
	"muse-admin/internal/define"
	knowledgeModel "muse-admin/internal/model/knowledge"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"
	"os"
	"strings"
	"time"
)

func AddExampleQuestion(f *types.ScheduleData) {
	exampleQuestion := types.ScheduleDataQuestion{
		QuestionNo:   cast.ToString(f.Id) + "_example",
		Type:         5,
		UsageType:    1, // 例题
		ReviewStatus: f.Status,
	}
	f.Question = append([]types.ScheduleDataQuestion{exampleQuestion}, f.Question...)
}

// convertMaterial 转换为kn_material表可录入的数据
func convertMaterial(data *types.ScheduleDataQuestion) *knowledgeModel.Material {
	return &knowledgeModel.Material{
		Title:  strings.TrimSpace(data.Material.Title),
		Author: strings.TrimSpace(data.Material.Author),
		Source: strings.TrimSpace(data.Material.Source),
		Content: sql.NullString{
			String: strings.TrimSpace(data.Material.Content),
			Valid:  true,
		},
		Background: sql.NullString{
			String: strings.TrimSpace(data.Material.Background),
			Valid:  true,
		},
		AuthorIntro: sql.NullString{
			String: strings.TrimSpace(data.Material.AuthorIntro),
			Valid:  true,
		},
	}
}

// convertExample 转换为kn_example表可录入的数据
func convertExample(data *types.ScheduleData, title, subTitle string) *knowledgeModel.Example {
	return &knowledgeModel.Example{
		Content:          "",                                      // 旧版
		Explain:          strings.TrimSpace(data.Example.Content), // 新版
		PlainText:        strings.TrimSpace(data.Example.PlainText),
		PlainTextArticle: strings.TrimSpace(data.Example.PlainTextArticle),
		Title:            strings.TrimSpace(title),
		SubTitle:         strings.TrimSpace(subTitle),
		LessonNotes:      strings.TrimSpace(data.Example.LessonNotes),
	}
}

// convertQuestion 转换为kn_question表可录入的数据
func convertQuestion(data *types.ScheduleDataQuestion, nodeType, materialId, exampleId int64) *knowledgeModel.LessonQuestion {
	return &knowledgeModel.LessonQuestion{
		QuestionNo:   data.QuestionNo,
		NodeType:     nodeType,
		Type:         cast.ToInt64(data.Type),
		UsageType:    cast.ToInt64(data.UsageType),
		GradePhase:   0,
		ReviewStatus: data.ReviewStatus,
		Level:        data.Level,
		MaterialId:   materialId,
		ExampleId:    exampleId,
		Ask:          strings.TrimSpace(data.Ask),
		Answer: sql.NullString{
			String: strings.TrimSpace(data.Answer),
			Valid:  true,
		},
		Analysis: sql.NullString{
			String: strings.TrimSpace(data.Analysis),
			Valid:  true,
		},
		StartTts: sql.NullString{
			String: strings.TrimSpace(data.StartTts),
			Valid:  true,
		},
		Source: strings.TrimSpace(data.Source),
	}
}

// convertOption 转换为kn_question_option表可录入的数据
func convertOption(data *types.ScheduleDataQuestion, questionId int64) []knowledgeModel.LessonQuestionOption {
	// 题库选项列表
	var optionList []knowledgeModel.LessonQuestionOption
	// 若有选项，需要处理为可以直接入目标表的选项数据
	if data.Type == SingleChoice {
		// 选项处理
		if len(data.Option) > 0 {
			for _, op := range data.Option {
				var isAnswer int64 // 是否正确,默认为0
				if data.Answer == op.Index {
					isAnswer = 1
				}
				optionList = append(optionList, knowledgeModel.LessonQuestionOption{
					LessonQuestionId: questionId,
					Sequence:         op.Sequence,
					OptionLabel:      op.Index,
					Content:          op.Value,
					IsAnswer:         isAnswer,
				})
			}
		}
	}
	return optionList
}

// convertLesson 转换为kn_lesson表可录入的数据
func convertLesson(data *types.ScheduleData, groupNo string) *knowledgeModel.Lesson {
	return &knowledgeModel.Lesson{
		LessonNo:      cast.ToString(data.Id),
		LessonGroupNo: groupNo,
		NodeType:      data.NodeType,
		ParentId:      data.ParentId,
		Level:         data.Level,
		Name:          data.LessonName,
		LessonType:    define.LessonTypeMapping[data.NodeType], // 小语文大语文统一在业务中称「主线课」
		ReviewStatus:  data.Status,
		Remark:        data.CourseName,
	}
}

// convertLessonResource 转换为kn_lesson_resource表可录入的数据
func convertLessonResource(lessonId, questionId int64, f *types.ScheduleData) *knowledgeModel.LessonResource {
	var serialNo int64
	// 小灶课时，有排序字段
	if f.NodeType == define.NodeType.SpecialCourse {
		// 若排序字段为0，直接取ID
		if f.Example.SerialNo == 0 {
			serialNo = f.Id
		} else {
			serialNo = f.Example.SerialNo
		}
	}

	return &knowledgeModel.LessonResource{
		LessonId:   lessonId,
		QuestionId: questionId,
		Sequence:   cast.ToFloat64(serialNo),
	}
}

// convertLessonPoint 转换为kn_lesson_point表可录入的数据
func convertLessonPoint(data *types.ScheduleData, lessonId int64) []knowledgeModel.LessonPoint {
	var point []knowledgeModel.LessonPoint

	// 按,拆分字符串
	parts := strings.Split(data.PointName, ",")

	// 输出拆分后的切片
	for _, part := range parts {
		point = append(point, knowledgeModel.LessonPoint{
			LessonId: lessonId,
			Name:     part,
		})
	}
	return point
}

// convertLessonQuestionTTS 转换为kn_lesson_question_tts表可录入的数据
func convertLessonQuestionTTS(ctx context.Context, cos *tencent.Cos, questionId int64, data string, miniMaxConf config.MiniMax, domain string) (*knowledgeModel.LessonQuestionTts, error) {
	if cos == nil {
		return nil, errs.NewMsg(errs.ServerErrorCode, "初始化腾讯云文件上传对象失败！")
	}

	audioReq := minimax.AudioReq{
		Text:   data,
		Model:  "speech-01-240228",
		Stream: false,
		TimberWeights: []minimax.TimberWeight{
			{
				VoiceId: "douyuxingchen_doulaoshi_1222",
				Weight:  1,
			},
		},
		VoiceSetting: minimax.VoiceSetting{
			VoiceId: "",
			Speed:   1.0,
			Vol:     2.0,
			Pitch:   0,
		},
		AudioSetting: minimax.AudioSetting{
			AudioSampleRate: 32000,
			Bitrate:         128000,
			Format:          "mp3",
		},
	}

	// 超时时间
	deadline, _ := ctx.Deadline()
	// 请求文本转语音
	minimaxVoice := minimax.NewVoice(miniMaxConf.ApiKey)
	resp, err := minimaxVoice.TextToAudio(ctx, miniMaxConf.GroupIds.CiYuan, audioReq, deadline.Sub(time.Now()))
	if err != nil {
		return nil, errs.WithMsg(err, errs.ServerErrorCode, "minimax语音接口异常")
	}

	// 音频编码写入临时文件
	tmpFilePath := "/tmp/audio/" + util.GetUUID() + ".mp3"
	if err = util.WriteHexCodeToFile(resp.Data.Audio, tmpFilePath); err != nil {
		return nil, errs.WithMsg(err, errs.ServerErrorCode, "音频编码写入临时文件失败")
	}
	defer func() {
		err = os.Remove(tmpFilePath)
		if err != nil {
			logc.Errorf(ctx, "删除临时文件失败，outputPath：%s，error：%s", tmpFilePath, err)
		}
	}()

	// 音频文件上传到腾讯云COS
	originFilePath := define.AuthOssPath.Audio + "/" + util.GetUUID()
	audioUrl, err := cos.Upload(tmpFilePath, originFilePath, domain)
	if err != nil {
		return nil, errs.WithMsg(err, errs.ServerErrorCode, "音频文件上传到腾讯云COS失败")
	}

	return &knowledgeModel.LessonQuestionTts{
		QuestionId: questionId,
		Content:    data,
		Url:        audioUrl,
		ApplyType:  1,
	}, nil
}
