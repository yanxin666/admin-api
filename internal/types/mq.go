package types

// ExcelImportBenefitModel 权益MQ模型
type ExcelImportBenefitModel struct {
	BenefitName  string `json:"benefit_name"`  // 权益名称
	Phone        string `json:"phone"`         // 手机号
	ChannelOrder string `json:"channel_order"` // 渠道订单号
	OrderNo      string `json:"order_no"`      // 订单号
	Source       int64  `json:"source"`        // 来源 1: 豆伴匠
	OrderStatus  int64  `json:"order_status"`  // 订单状态
}

// ExcelImportBalanceData 用户余额MQ模型
type ExcelImportBalanceData struct {
	Phone   string `json:"phone"`     // 用户手机号
	Amount  int64  `json:"amount"`    // 金额
	OrderNo string `json:"source_no"` // 订单号
}

// ImportData 导入MQ模型
type ImportData struct {
	Id int64 `json:"id"`
}

// ScheduleData 规划MQ模型
type ScheduleData struct {
	Source     int64  `json:"source"`      // 数据来源方 1.九霄-阅读理解课堂
	Id         int64  `json:"id"`          // ID
	NodeType   int64  `json:"node_type"`   // 类型 1.小灶课[原理讲解] 2.小语文[目标课] 3.大语文[目标课]
	ParentId   int64  `json:"parent_id"`   // 只有大语文[目标课]才会使用，用作关联父子级
	Level      int64  `json:"level"`       // 只有小灶课才会使用,1-4个等级
	PointName  string `json:"point_name"`  // 知识点名称 小灶课为1个; 目标课为多个用“,”分割
	CourseName string `json:"course_name"` // 课程名称
	Status     int64  `json:"status"`      // 审核状态 1.初始录入 2.审核中 3.审核通过 4.审核未通过 5.下架
	Grade      int64  `json:"grade"`       // 年级
	Example    struct {
		Content          string `json:"content"`            // 例题内容，JSON格式
		PlainText        string `json:"plain_text"`         // JSON中的纯文本内容
		PlainTextArticle string `json:"plain_text_article"` // JSON中的纯文本正文
		Title            string `json:"title"`              // 标题
		SubTitle         string `json:"sub_title"`          // 副标题
		SerialNo         int64  `json:"serial_no"`          // 排序
		LessonNotes      string `json:"lesson_notes"`       // 笔记
	} `json:"example"` // 例题
	Question []ScheduleDataQuestion `json:"question"`

	// Deprecated: 暂未使用
	Term int64 `json:"term"` // 学期
	// Deprecated: 暂未使用
	UnitName string `json:"unit_name"` // 单元名称
	// Deprecated: 使用到了，但未做业务渲染，属于废物字段
	LessonName string `json:"lesson_name"` // 课节名称
}
type ScheduleDataQuestion struct {
	Material struct {
		Title       string `json:"title"`        // 素材名称
		Author      string `json:"author"`       // 作者
		Source      string `json:"source"`       // 素材来源
		Content     string `json:"content"`      // 素材内容，比如选取的阅读素材片段
		Background  string `json:"background"`   // 写作背景
		AuthorIntro string `json:"author_intro"` // 作者介绍
	} `json:"material"` // 素材
	QuestionNo   string                       `json:"question_no"`   // 题目编号
	Type         int64                        `json:"type"`          // 试题类型 1.单选 2.多选 3.填空 4.判断 5.简答 6.阅读题 7.作文
	UsageType    int64                        `json:"usage_type"`    // 使用类型 1.例题 2.练习题 3.候补题
	GradePhase   int64                        `json:"grade_phase"`   // 年级学段 1.小学 2.初中 3.高中
	ReviewStatus int64                        `json:"review_status"` // 审核状态 1.初始录入 2.审核中 3.审核通过 4.审核未通过 5.下架
	Level        int64                        `json:"level"`         // 等级
	Ask          string                       `json:"ask"`           // 问题
	Answer       string                       `json:"answer"`        // 答案
	Option       []ScheduleDataQuestionOption `json:"option"`        // 选项
	Analysis     string                       `json:"analysis"`      // 题目解析
	StartTts     string                       `json:"start_tts"`     // 开场白TTS
	Duration     int64                        `json:"duration"`      // 预计用时
	Source       string                       `json:"source"`        // 来源
	Version      int64                        `json:"version"`       // 版本
}
type ScheduleDataQuestionOption struct {
	Sequence int64  `json:"sequence"` // 选项展示顺序[正序]
	Index    string `json:"index"`    // 选项
	Value    string `json:"value"`    // 值
}

// WritePPTData 写作参谋PPT的MQ模型
type WritePPTData struct {
	Source            int64    `json:"source"`              // 数据来源方 2.九霄-写作参谋课堂
	Id                int64    `json:"id"`                  // ID
	Topic             string   `json:"topic"`               // 主题
	Unit              int64    `json:"unit"`                // 单元
	Series            string   `json:"series"`              // 大纲系列 例如：R2、R3等
	Title             string   `json:"title"`               // 标题
	LessonNumber      int64    `json:"lesson_number"`       // 课节序号
	LessonType        int64    `json:"lesson_type"`         // 课节类型 1.技巧 2.赏析
	LessonCategory    int64    `json:"lesson_category"`     // 课节分类 1.细节 2.表达 3.布局
	LessonContent     string   `json:"lesson_content"`      // 讲稿内容 JSON格式
	LessonContentText string   `json:"lesson_content_text"` // 讲稿内容 纯文本格式
	Notes             string   `json:"notes"`               // 笔记
	ReviewStatus      int64    `json:"review_status"`       // 审核状态 1.初始录入 2.审核中 3.审核通过 4.审核未通过 5.下架
	PPTFiles          []string `json:"ppt_files"`           // PPT文件URL JSON列表
	MindSlices        string   `json:"mind_slices"`         // 脑图步骤 JSON列表
	MindFull          string   `json:"mind_full"`           // 脑图全量
	AudioUrl          string   `json:"audio_url"`           // 讲稿音频地址
}

// LiveData 直播MQ模型
type LiveData struct {
	Source           int64      `json:"source"`             // 数据来源方 3.九霄-直播
	EnvType          int64      `json:"env_type"`           // 环境类型 1.预览(对标测试环境，课堂数据直接覆盖) 2.发布(对标生产环境，课堂数据为审核状态)
	ModeType         int64      `json:"mode_type"`          // 模式类型 1.模版模式 2.定制工厂(若为工厂模式，老师名称、音色、音频、举手等字段都需要有，若没有就报错)
	LiveNo           string     `json:"version_no"`         // 课程编号
	Name             string     `json:"name"`               // 课程名称
	Description      string     `json:"description"`        // 课程描述
	FullImage        string     `json:"course_bgimg"`       // 全屏背景URL
	TeacherCode      string     `json:"teacher_code"`       // 老师编码（唯一标识）
	TeacherName      string     `json:"teacher_name"`       // 老师名称
	TeacherNickname  string     `json:"teacher_nickname"`   // 老师昵称
	TeacherTone      string     `json:"teacher_tone"`       // 老师音色
	HandAudio        string     `json:"hand_audio"`         // 老师举手发声音频
	TeacherToneModel string     `json:"teacher_tone_model"` // 老师模型
	TeacherToneType  int64      `json:"teacher_tone_type"`  // 老师音色类型 1.minimax 2.火山
	TeacherIntroduce string     `json:"teacher_introduce"`  // 老师介绍
	Topic            string     `json:"topic"`              // 课堂开始时的寒暄主题
	IntervalDuration []int      `json:"interval_duration"`  // 寒暄话题周期(s)
	TeacherGif       string     `json:"avatar_gif"`         // 讲师头像动图
	TeacherPng       string     `json:"avatar_png"`         // 讲师头像静图
	TeacherVideo     string     `json:"primary"`            // 老师主视频
	Duration         int64      `json:"duration"`           // 老师主视频总时长
	PPTVideo         string     `json:"sub"`                // 副视频
	Draft            string     `json:"subtitle"`           // 字幕
	DraftRag         string     `json:"rag"`                // 纯文字字幕
	SimpleDraft      string     `json:"outline"`            // 提炼出来的精简字幕
	BlockTime        string     `json:"air_event"`          // 停顿气口
	Content          string     `json:"orchestrate"`        // 内容
	ReviewStatus     int64      `json:"review_status"`      // 审核状态 1.初始录入 2.审核中 3.审核通过 4.审核未通过 5.下架
	Roles            []Roles    `json:"roles"`              // 角色
	Navigate         []Navigate `json:"nav"`                // 导航进度条
	Note             string     `json:"note"`               // 课堂笔记
	Homework         string     `json:"homework"`           // 课后作业
	Account          string     `json:"account"`            // 推送操作人
}
type Roles struct {
	Character string `json:"character"`
	Role      string `json:"role"`
	UserId    int    `json:"user_id"`
	UserName  string `json:"user_name"`
}
type Navigate struct {
	Name      string  `json:"name"`       // 环节名称
	StartTime float64 `json:"start_time"` // 触发时间(s)
}

// 超能训练MQ模型
type (
	SuperTrainData struct {
		// 公共字段
		Source  int64  `json:"source,omitzero"`   // 数据来源方 4.九霄-超能训练
		EnvType int64  `json:"env_type,omitzero"` // 环境类型 1.预览(对标测试环境，课堂数据直接覆盖) 2.发布(对标生产环境，课堂数据为审核状态)
		Account string `json:"account,omitzero"`  // 推送操作人
		Status  int64  `json:"status,omitzero"`   // 审核状态 1.初始录入 2.审核中 3.审核通过 4.审核未通过 5.下架

		// 课程信息
		CourseType  int64       `json:"course_type,omitzero"`  // 课程类型 1.超练作文 2.超练阅读
		CourseNo    string      `json:"course_no,omitzero"`    // 课程编号
		CourseName  string      `json:"course_name,omitzero"`  // 课程名称
		LessonName  string      `json:"lesson_name,omitzero"`  // 课节名称
		Subject     int64       `json:"subject,omitzero"`      // 学科 1.语文
		Unit        string      `json:"unit,omitzero"`         // 单元名
		Image       string      `json:"image,omitzero"`        // 课程封面图
		Intro       string      `json:"intro,omitzero"`        // 课程介绍
		OpenTime    string      `json:"open_time,omitzero"`    // 开课时间，暂时不用，后续可能会用到
		Type        int64       `json:"type,omitzero"`         // 类型 0.免费课 1.试听课 2.付费课 暂时不用，后续可能会用到
		CourseExtra CourseExtra `json:"course_extra,omitzero"` // 课程额外字段
		Chapter     []*Chapter  `json:"chapter,omitzero"`      // 章节信息：每个课程下有多个章节
	}
	// Chapter 章节信息
	Chapter struct {
		No         string       `json:"no,omitzero"`            // 章节编号
		Name       string       `json:"name,omitzero"`          // 章节标题
		Title      string       `json:"title,omitzero"`         // 章节精简版标题
		Image      string       `json:"image,omitzero"`         // 章节封面图
		Intro      string       `json:"intro,omitzero"`         // 章节介绍
		Type       int64        `json:"type,omitzero"`          // 章节类型 1.任务 2.文章
		Sequence   float64      `json:"sequence,omitzero"`      // 章节排序
		GuideVideo string       `json:"guide_video,omitzero"`   // 引导视频URL
		TeacherId  int64        `json:"teacher_id,omitzero"`    // 老师ID
		Extra      ChapterExtra `json:"chapter_extra,omitzero"` // 章节额外字段
		Task       []*Task      `json:"task,omitzero"`          // 任务信息：每个章节下有多个任务
	}
	// Task 任务信息
	Task struct {
		TaskNo         string     `json:"task_no,omitzero"`         // 任务编号
		Title          string     `json:"title,omitzero"`           // 标题
		Subtitle       string     `json:"subtitle,omitzero"`        // 子标题
		Description    string     `json:"description,omitzero"`     // 描述
		ProgressPrefix string     `json:"progress_prefix,omitzero"` // 进度前缀、时空稳定值、时空通道已建立等
		Type           int64      `json:"type,omitzero"`            // 类型 1.文章段落 2.身份卡 3.文章总结 4.时空长廊 200.主线任务 201.简单副本任务 202 困难副本
		Sequence       float64    `json:"sequence,omitzero"`        // 排序
		Extra          TaskExtra  `json:"task_extra,omitzero"`      // 任务额外字段
		SubTask        []*SubTask `json:"sub_task,omitzero"`        // 子任务信息：每个任务下有多个子任务
	}
	// SubTask 子任务信息
	SubTask struct {
		SubTaskNo   string       `json:"sub_task_no,omitzero"` // 子任务编号
		Mode        int64        `json:"mode,omitzero"`        // 子任务模式 1.简单模式 2.困难模式 (当 SuperTrainData.CourseType = 1时，默认传1)
		Type        int64        `json:"type,omitzero"`        // 类型 1.普通任务 2.身份卡 3.聊天 4.结尾字幕 5.签名 6.时空长廊
		Title       string       `json:"title,omitzero"`       // 标题
		Description string       `json:"description,omitzero"` // 描述
		ImageUrl    string       `json:"image_url,omitzero"`   // 图片地址
		Sequence    float64      `json:"sequence,omitzero"`    // 排序
		Extra       SubTaskExtra `json:"sub_extra,omitzero"`   // 子任务额外字段

		// 当 SuperTrainData.CourseType = 2 时有值
		Article string `json:"article,omitzero"` // 原文资料
		Medal   Medal  `json:"medal,omitzero"`   // 每个子任务都对应一个奖牌配置
	}
)

// ==================== 额外信息结构 ====================
type (
	// CourseExtra 课程额外信息
	CourseExtra struct {
		// --- 课程描述相关 --- 仅阅读使用
		Title              string `json:"title,omitzero"`                 // 首页课程标题
		ButtonText         string `json:"button_text,omitzero"`           // 首页按钮文字
		TeacherName        string `json:"teacher_name,omitzero"`          // 首页老师名字
		TeacherVideo       string `json:"teacher_video,omitzero"`         // 首页老师视频m3u8
		TeacherVideoWeb    string `json:"teacher_video_web,omitzero"`     // 首页老师视频mp4
		HomePageGuideVideo string `json:"home_page_guide_video,omitzero"` // 开场视频地址

		// --- 课程背景 --- 仅写作使用
		CourseBackground  string `json:"course_background,omitempty"`  // 课程背景描述
		CourseAppellation string `json:"course_appellation,omitempty"` // 课程称谓描述

		// --- 最后一关成文 --- 仅写作使用
		Through           string   `json:"through"`             // 穿越名
		ArticleTts        string   `json:"article_tts"`         // 作文选字数tts
		SecretsTask       string   `json:"secrets_task"`        // 探秘任务名
		ArticleTitle      string   `json:"article_title"`       // 作文题目
		OutlineTTSText    []string `json:"outline_tts_text"`    // 写作大纲tts文本
		AfterArticleTts   string   `json:"after_article_tts"`   // 生成作文后的tts
		WritingSkills     string   `json:"writing_skills"`      // 写作技法 大模型使用
		ArticleCenterIdea string   `json:"article_center_idea"` // 作文中心思想 大模型使用
	}
	// ChapterExtra 章节额外信息
	ChapterExtra struct {
		// --- 公用字段  --- 写作与阅读都需要以下字段
		GuideVideo     string `json:"guide_video,omitzero"`      // 引导页视频 Url
		GuideVideoTime int    `json:"guide_video_time,omitzero"` // 引导视频时长

		// --- 引导关 --- 仅写作使用
		GuideTask         []CourseChapterGuideTask `json:"guide_task,omitempty"`         // 引导关数据
		WritingTechniques string                   `json:"writing_techniques,omitempty"` // 写作技法
		CompositionTitle  string                   `json:"composition_title,omitempty"`  // 作文标题

		// --- 章节基本信息 (@王俊) --- 仅写作使用
		WritingGoals      string `json:"writing_goals,omitzero"`        // 写作目标
		FullStarAdvice    string `json:"full_star_advice,omitzero"`     // 文章页面满星建议
		ExcellentSample   string `json:"excellent_sample,omitzero"`     // 优秀范文示例
		ArticleTopTipText string `json:"article_top_tip_text,omitzero"` // 文章页面顶部提示文本

		// --- 引导信息 (@张兴、@梁立凯) --- 仅写作使用
		BackgroundUrl  string `json:"background_url,omitzero"`   // 引导背景图
		GuideVideoDesc string `json:"guide_video_desc,omitzero"` // 引导页视频文案
		GuideVideoUrl  string `json:"guide_video_url,omitzero"`  // 引导页视频 Url 需要与上方 GuideVideo 保持一致，这里代码中会手动处理，数据组不需要传参

		Question []CourseChapterQuestion `json:"question,omitempty"` // 引导关问题列表
	}
	// CourseChapterQuestion 章节引导关问题
	CourseChapterQuestion struct {
		Clue       string   `json:"clue,omitempty"`       // 问题内容
		Option     []string `json:"option,omitempty"`     // 选项列表
		Topic      string   `json:"topic,omitempty"`      // 题目
		Image      string   `json:"image,omitempty"`      // 图片地址
		Sequence   int64    `json:"sequence,omitempty"`   // 序号
		Background string   `json:"background,omitempty"` // 背景
		Key        string   `json:"key,omitempty"`        // 答案及解析
	}
	// TaskExtra 任务额外信息
	TaskExtra struct {
		// 是否按年级的达标字数来计分，若为true，则下方的 SubTaskExtra.GradeRules 必传
		ParagraphType bool `json:"paragraph_type,omitzero"` // 文章段落计分类型：1.按年级的达标字数来计分 2.按任务来计分

		// 如果该字段不为空，则表示该任务下有个别子任务的 ShowRule 字段里存在指定段落写作，需要固定格式：${type:chapter_user_article|chapter_no=Chapter.No}
		ParagraphAppoint string `json:"paragraph_appoint,omitzero"` // 指定段落编号, Chapter.No

		// Deprecated: 以下格式都于2025年10月30日全都已废弃
		// // --- 任务额外信息 (@李斌) --- 仅写作使用
		// ProgressLT100 string           `json:"progress_lt_100,omitzero"` // 进度小于100%时的标题
		// ProgressEQ100 string           `json:"progress_eq_100,omitzero"` // 进度等于100%时的提示文案
		// BackgroundUrl string           `json:"background_url,omitzero"`  // 背景图
		// Guide         []TaskExtraGuide `json:"guide,omitzero"`           // 引导信息
	}
	// SubTaskExtra 子任务额外信息
	SubTaskExtra struct {
		// --- 公用字段  --- 写作与阅读都需要以下字段
		Question string `json:"question,omitzero"` // 超练写作题目 / 超练阅读问题

		// --- 超练阅读配置 (@宋听森) --- 仅阅读使用
		Answer       string `json:"answer,omitzero"`        // 问题答案
		ScorePoint   string `json:"score_point,omitzero"`   // 得分点
		AnswerMethod string `json:"answer_method,omitzero"` // 答题方式
		QuestionType int64  `json:"question_type,omitzero"` // 题目类型 0.非选择题 1.选择题

		// --- 超练阅读配置 (@李超) --- 仅阅读使用
		Tts           string `json:"tts,omitzero"`             // tts
		TeacherTts    string `json:"teacher_tts,omitzero"`     // 教师tts
		StoryData     string `json:"story_data,omitzero"`      // 情景导入
		ArticleMp3    string `json:"article_mp3,omitzero"`     // 文章mp3
		TtsMp3        string `json:"tts_mp3,omitzero" `        // tts mp3
		TeacherTtsMp3 string `json:"teacher_tts_mp3,omitzero"` // 教师tts mp3
		ArticleDraft  string `json:"article_draft,omitzero"`   // 文章断句
		Ditto         bool   `json:"ditto,omitzero"`           // 是否原文和剧情同上

		// 以下都为写作穿越的数据
		// --- 模型参数配置 (@王鹏飞) --- 仅写作使用
		ExemplarySentences []string `json:"exemplary_sentences,omitzero"` // 优秀句子示例
		AnswerKey          []string `json:"answer_key,omitzero"`          // 答案关键字
		Goals              string   `json:"goals,omitzero"`               // 写作要求
		Storyline          string   `json:"storyline,omitzero"`           // 故事线索
		Approach           string   `json:"approach,omitzero"`            // 写作思路
		ShowRule           ShowRule `json:"show_rule,omitzero"`           // 展示规则
		// 不同年级达标字数计算规则，默认可不填，当 TaskExtra.ParagraphType = true 时必传
		GradeRules []GradeRules `json:"grade_rules,omitzero"`

		// Deprecated: 以下格式都于2025年10月30日全都已废弃
		// // --- 展示线索配置 (@李斌) --- 仅写作使用
		// ScanImageListUrl []ScanImageListUrl `json:"scan_image_list_url,omitzero"` // 线索扫描图列表地址
		// MyImageUrl       string             `json:"my_image_url,omitzero"`        // 我的线索图片地址
		// IsAuto           bool               `json:"is_auto,omitzero"`             // 是否自动执行
		//
		// // --- 学生卡片信息配置 (@赵雨龙) --- 仅写作使用
		// StudentCard StudentCard  `json:"student_card,omitzero"` // 学生卡片信息
		// GradeRules  []GradeRules `json:"grade_rules,omitzero"`  // 学生卡片信息
		//
		// // --- 线索详情图片地址 (@张振彪) --- 仅写作使用
		// ImageUrl string `json:"image_url,omitzero"` // 线索详情图片地址
		//
		// // --- 签名墙 (@张兴) --- 仅写作使用
		// SignatureData SignatureData `json:"signature_data,omitzero"` // 签名墙数据
		//
		// // --- 结尾字幕信息 (@王俊) --- 仅写作使用
		// EndingCaption EndingCaption `json:"ending_caption,omitzero"` // 结尾字幕信息
	}

	// CourseChapterGuideTask 章节引导关数据结构
	CourseChapterGuideTask struct {
		VideoUrl    string                        `json:"video_url,omitempty"`     // 视频地址 (m3u8)
		VideoWebUrl string                        `json:"video_web_url,omitempty"` // 视频地址 (mp4)
		Content     string                        `json:"content,omitempty"`       // 内容
		Images      []CourseChapterGuideTaskImage `json:"images,omitempty"`        // 图片列表
	}

	// CourseChapterGuideTaskImage 章节引导关图片数据
	CourseChapterGuideTaskImage struct {
		Url     string                             `json:"url,omitempty"`     // 图片地址
		Theme   string                             `json:"theme,omitempty"`   // 主题
		Title   string                             `json:"title,omitempty"`   // 标题
		Trigger CourseChapterGuideTaskImageTrigger `json:"trigger,omitempty"` // 触发事件
		Time    int64                              `json:"time,omitempty"`    // 时间
	}

	// CourseChapterGuideTaskImageTrigger 章节引导关图片触发事件
	CourseChapterGuideTaskImageTrigger struct {
		Type    int64  `json:"type,omitempty"`    // 事件类型
		Payload string `json:"payload,omitempty"` // 事件数据 (JSON字符串)
	}

	// TaskExtraGuide 任务引导信息配置
	TaskExtraGuide struct {
		Role          int64  `json:"role,omitzero"`           // 角色类型 1.老师 2.NPC 3.旁白 4.手势引导
		RoleVoiceKey  string `json:"role_voice_key,omitzero"` // 角色语音key
		Name          string `json:"name,omitzero"`           // 名称
		Url           string `json:"url,omitzero"`            // 资源地址
		UrlType       int64  `json:"url_type,omitzero"`       // 资源类型 1.图片 2.视频 3.音频 4.链接 5.文本
		Content       string `json:"content,omitzero"`        // 文案内容
		BackgroundUrl string `json:"background_url,omitzero"` // 背景图
		WebData       string `json:"web_data,omitzero"`       // 前端数据
	}

	// ScanImageListUrl 子任务图片资源配置
	ScanImageListUrl struct {
		Title   string `json:"title,omitzero"`    // 子任务标题
		Url     string `json:"url,omitzero"`      // 子任务地址
		WebData string `json:"web_data,omitzero"` // 前端数据
	}

	// StudentCard 学生卡片信息
	StudentCard struct {
		Age     string `json:"age,omitzero"`     // 年龄
		Address string `json:"address,omitzero"` // 地址
		School  string `json:"school,omitzero"`  // 学校
	}

	// ArticleDraft 文章断句信息
	ArticleDraft struct {
		Text      string `json:"text,omitzero"`       // 断句内容
		TimeBegin int    `json:"time_begin,omitzero"` // 断句开始时间
		TimeEnd   int    `json:"time_end,omitzero"`   // 断句结束时间
	}

	// Medal 奖牌配置(阅读理解使用)
	Medal struct {
		Name            string `json:"name,omitzero"`             // 名称
		Description     string `json:"description,omitzero"`      // 描述
		Type            int64  `json:"type,omitzero"`             // 类型 1.金钥匙 2.银钥匙 3.铜钥匙
		PromptKnowledge string `json:"prompt_knowledge,omitzero"` // p工程类关联知识点
	}

	// GradeRules 根据用户年级判断年级范围内匹配的达标字数
	GradeRules struct {
		Min   int `json:"min,omitzero"`   // 最小年级
		Max   int `json:"max,omitzero"`   // 最大年级
		Chars int `json:"chars,omitzero"` // 需要达标的字数
	}

	// EndingCaption 结尾字幕信息
	EndingCaption struct {
		Text     string `json:"text,omitzero"`      // 字幕文本
		AudioUrl string `json:"audio_url,omitzero"` // 音频地址
	}

	// SignatureData 签名墙数据
	SignatureData struct {
		SignatureText      string `json:"signature_text,omitzero"`       // 签名墙文本
		GuideVideoUrl      string `json:"guide_video_url,omitzero"`      // 引导视频链接
		GuideVideoTime     int64  `json:"guide_video_time,omitzero"`     // 引导视频时间,单位秒
		GuideVideoCover    string `json:"guide_video_cover,omitzero"`    // 引导视频封面图
		GuideText          string `json:"guide_text,omitzero"`           // 引导文本
		GuideBackgroundUrl string `json:"guide_background_url,omitzero"` // 引导背景图
		EndTtsText         string `json:"end_tts_text,omitzero"`         // 结束TTS文本
	}

	// ShowRule 端上展示时，所涉及到的所有规则
	ShowRule struct {
		IsQuestion bool      `json:"is_question,omitzero"` // 是否有题目
		Intro      []Content `json:"intro,omitzero"`       // 介绍
		Settlement []Content `json:"settlement,omitzero"`  // 结算
	}
	Content struct {
		// 类型
		// 1.背景图+视频+文本
		// 2.背景图+点位+TTS文本
		// 3.背景图+视频+视频文本+锚点
		// 4.背景图+视频+视频文本+叠加图+锚点
		// 5.计时切图+视频+视频文本+图片+时间列表
		Type    int64  `json:"type,omitzero"`
		Payload string `json:"payload,omitzero"` // json序列化字符串，格式由前端定义
	}
)

//
// type SubTaskExtra struct {
// 	ShowRule ShowRule `json:"showRule,omitzero"` // 展示规则
// 	Write    Write    `json:"write,omitzero"`    // 超练写作
// 	Read     Read     `json:"read,omitzero"`     // 超练阅读
// }
//
// type (
// 	// ShowRule 端上展示时，所涉及到的所有规则
// 	ShowRule struct {
// 		IsAuto     bool       `json:"isAuto,omitzero"`     // 是否自动执行
// 		IsQuestion bool       `json:"question,omitzero"`   // 是否有题目
// 		Intro      []Intro    `json:"intro,omitzero"`      // 介绍
// 		Settlement Settlement `json:"settlement,omitzero"` // 结算
// 	}
// 	Intro struct {
// 		Type    int     `json:"type,omitzero"`    // 类型 1.背景图+视频+文本 2.背景图+点位+TTS文本
// 		Payload Payload `json:"payload,omitzero"` // json序列化字符串，格式由前端定义
// 	}
// 	Payload struct {
// 		BackgroundUrl string `json:"backgroundUrl,omitzero"` // 背景图
// 		Tts           string `json:"tts,omitzero"`           // TTS文本
// 	}
// 	Settlement struct {
// 		Background string `json:"background,omitzero"` // 背景图
// 		Tts        string `json:"tts,omitzero"`        // TTS文本
// 	}
// )
//
// // 写作穿越使用
// type (
// 	Write struct {
// 		Question       string         `json:"question,omitzero"`       // 题目
// 		ImageUrl       string         `json:"imageUrl,omitzero"`       // 线索详情图片地址
// 		DialoguePrompt DialoguePrompt `json:"dialoguePrompt,omitzero"` // 对话所需要的p工程
// 	}
// 	// DialoguePrompt 对话所需要的p工程
// 	DialoguePrompt struct {
// 		ExemplarySentences string `json:"exemplarySentences,omitzero"` // 优秀句子示例
// 		AnswerKey          string `json:"answerKey,omitzero"`          // 答案关键字
// 		Goals              string `json:"goals,omitzero"`              // 写作要求
// 		Storyline          string `json:"storyline,omitzero"`          // 故事线索
// 		Approach           string `json:"approach,omitzero"`           // 写作思路
// 	}
// )
//
// // 阅读理解使用
// type (
// 	Read struct {
// 		Prepare    Prepare    `json:"prepare,omitzero"`    // 课程互动
// 		SceneStory SceneStory `json:"sceneStory,omitzero"` // 情景故事
// 		Question   Question   `json:"question,omitzero"`   // 题目
// 		Medal      Medal      `json:"medal,omitzero"`      // 奖牌配置
// 	}
// 	// Prepare 课程互动
// 	Prepare struct {
// 		ArticleDraft  string `json:"articleDraft,omitzero"`
// 		Tts           string `json:"tts,omitzero"`
// 		TtsMp3        string `json:"ttsMp3,omitzero"`
// 		TeacherTts    string `json:"teacherTts,omitzero"`
// 		TeacherTtsMp3 string `json:"teacherTtsMp3,omitzero"`
// 	}
// 	// SceneStory 情景故事
// 	SceneStory struct {
// 		Data  string `json:"data,omitzero"`  // 内容
// 		Ditto bool   `json:"ditto,omitzero"` // 是否原文和剧情同上
// 	}
// 	// Question 题目
// 	Question struct {
// 		Article      string `json:"article,omitzero"`      // 原文资料
// 		ArticleMp3   string `json:"articleMp3,omitzero"`   // 文章mp3
// 		QuestionType string `json:"questionType,omitzero"` // 题目类型 0.非选择题 1.选择题
// 		Question     string `json:"question,omitzero"`     // 题目
// 		Answer       string `json:"answer,omitzero"`       // 问题答案
// 		ScorePoint   string `json:"scorePoint,omitzero"`   // 得分点
// 		AnswerMethod string `json:"answerMethod,omitzero"` // 答题方式
// 	}
// 	// Medal 奖牌配置
// 	Medal struct {
// 		Name            string `json:"name,omitzero"`             // 名称
// 		Description     string `json:"description,omitzero"`      // 描述
// 		Type            int64  `json:"type,omitzero"`             // 类型 1.金钥匙 2.银钥匙 3.铜钥匙
// 		PromptKnowledge string `json:"prompt_knowledge,omitzero"` // p工程类关联知识点
// 	}
// )
