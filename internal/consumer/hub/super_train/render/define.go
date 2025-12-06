package render

type (
	// CourseExtra 课程额外信息
	CourseExtra struct {
		// --- 课程信息相关 (@张金涛) ---
		Title           string `json:"title,omitzero"`             // 首页课程标题
		SubTitle        string `json:"subTitle,omitzero"`          // 首页副标题
		TeacherName     string `json:"teacherName,omitzero"`       // 首页老师名字
		TeacherVideo    string `json:"teacherVideo,omitzero"`      // 首页老师视频 m3u8
		TeacherVideoWeb string `json:"teacher_video_web,omitzero"` // 首页老师视频mp4
		ButtonText      string `json:"buttonText,omitzero"`        // 首页按钮文字
		VideoTime       int64  `json:"videoTime,omitzero"`         // 视频时长

		// --- 课程描述相关 (@严鑫) ---
		CourseForm       string `json:"courseForm,omitzero"`       // 课程形式描述
		CourseStory      string `json:"courseStory,omitzero"`      // 课程故事描述
		RolePlaying      string `json:"rolePlaying,omitzero"`      // 角色扮演描述
		CourseObjectives string `json:"courseObjectives,omitzero"` // 课程目标描述
	}
	// ChapterExtra 章节额外信息
	ChapterExtra struct {
		// --- 章节基本信息 (@王俊) ---
		ArticleBgUrl         string `json:"articleBgUrl,omitzero"`         // 文章背景图
		WritingGoals         string `json:"writingGoals,omitzero"`         // 写作目标
		FullStarAdvice       string `json:"fullStarAdvice,omitzero"`       // 文章页面满星建议
		ExcellentSample      string `json:"excellentSample,omitzero"`      // 优秀范文示例
		ArticleTopTipText    string `json:"articleTopTipText,omitzero"`    // 文章页面顶部提示文本
		ArticleEndAttachText string `json:"articleEndAttachText,omitzero"` // 文章结尾附加文本
		HasEndingTask        int    `json:"hasEndingTask,omitzero"`        // 是否有结束任务(如字幕签名) 0.没有 1.有

		// --- 引导信息 (@张兴) ---
		GuideVideo     string `json:"guideVideo,omitzero"`     // 引导页视频 Url
		BackgroundUrl  string `json:"backgroundUrl,omitzero"`  // 引导背景图
		GuideVideoTime int    `json:"guideVideoTime,omitzero"` // 引导视频时长
		GuideContent   string `json:"guideContent,omitzero"`   // 引导文本

		// --- 引导模型参数配置 (@梁立凯) ---
		GuideVideoUrl  string `json:"guideVideoUrl,omitzero"`  // 引导页视频 Url 需要与上方 GuideVideo 保持一致
		GuideVideoDesc string `json:"guideVideoDesc,omitzero"` // 引导页视频文案
	}
	// TaskExtra 任务额外信息
	TaskExtra struct {
		ParagraphType int64 `json:"paragraphType,omitempty"` // 文章段落计分类型：0或1.按字数计分 2.按任务计分
	}
	// SubTaskExtra 子任务额外信息
	SubTaskExtra struct {
		// --- 公用字段  ---
		Question string `json:"question,omitzero"` // 超练写作题目 / 超练阅读问题

		// --- 超练阅读配置 (@宋听森) ---
		Answer       string `json:"answer,omitzero"`       // 问题答案
		ScorePoint   string `json:"scorePoint,omitzero"`   // 得分点
		AnswerMethod string `json:"answerMethod,omitzero"` // 答题方式
		QuestionType int64  `json:"questionType,omitzero"` // 题目类型 0.非选择题 1.选择题

		// --- 超练阅读配置 (@李超) ---
		Tts           string         `json:"tts,omitzero"`           // tts
		TeacherTts    string         `json:"teacherTts,omitzero"`    // 教师tts
		StoryData     string         `json:"storyData,omitzero"`     // 情景导入
		ArticleMp3    string         `json:"articleMp3,omitzero"`    // 文章mp3
		TtsMp3        string         `json:"ttsMp3,omitzero" `       // tts mp3
		TeacherTtsMp3 string         `json:"teacherTtsMp3,omitzero"` // 教师tts mp3
		ArticleDraft  []ArticleDraft `json:"articleDraft,omitzero"`  // 文章断句

		// 以下都为写作穿越的数据
		// --- 模型参数配置 (@王鹏飞) ---
		ExemplarySentences []string     `json:"exemplarySentences,omitzero"` // 优秀句子示例
		AnswerKey          []string     `json:"answerKey,omitzero"`          // 答案关键字
		Goals              string       `json:"goals,omitzero"`              // 写作要求
		Storyline          string       `json:"storyline,omitzero"`          // 故事线索
		Approach           string       `json:"approach,omitzero"`           // 写作思路
		GradeRules         []GradeRules `json:"gradeRules,omitzero"`         // 学生卡片信息
	}

	// ArticleDraft 文章断句信息
	ArticleDraft struct {
		Text      string `json:"text,omitzero"`      // 断句内容
		TimeBegin int    `json:"timeBegin,omitzero"` // 断句开始时间
		TimeEnd   int    `json:"timeEnd,omitzero"`   // 断句结束时间
	}

	// Medal 奖牌配置(阅读理解使用)
	Medal struct {
		Name            string `json:"name,omitzero"`            // 名称
		Description     string `json:"description,omitzero"`     // 描述
		Type            int64  `json:"type,omitzero"`            // 类型 1.金钥匙 2.银钥匙 3.铜钥匙
		PromptKnowledge string `json:"promptKnowledge,omitzero"` // p工程类关联知识点
	}

	// GradeRules 根据用户年级判断年级范围内匹配的达标字数
	GradeRules struct {
		Min   int `json:"min,omitzero"`   // 最小年级
		Max   int `json:"max,omitzero"`   // 最大年级
		Chars int `json:"chars,omitzero"` // 需要达标的字数
	}
)
