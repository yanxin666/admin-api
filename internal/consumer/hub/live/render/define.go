package render

// ContentDraft 大内容字幕
type ContentDraft struct {
	Text      string `json:"text"`       // 字幕内容
	TimeBegin int64  `json:"time_begin"` // 开始时间
	TimeEnd   int64  `json:"time_end"`   // 结束时间
}

// Content 大内容解析
type Content struct {
	Links []Link `json:"links"`
}

// Link - 每个链接的结构体
type Link struct {
	Type     string     `json:"type"`     // "SMALL_CHAT":班主任闲聊 | "TALK":寒暄 | "IN_CLASS":上课
	Timeline []Timeline `json:"timeline"` // 时间线
	Screen   struct {
		Image struct {
			Id    int64  `json:"id"`    // ID
			Media string `json:"media"` // URL
		} `json:"image"` // 视频开场背景图
	} `json:"screen,omitempty"` // 选项列表
}

type Screen struct {
	Image string `json:"image"` // 视频开场背景URL
}

// Timeline - 事件时间轴信息
type Timeline struct {
	// "user_msg":用户信息
	// "turn_on_camera":打开摄像头
	// "game_select":事件互动多选
	// "game_input":事件互动单选
	// "homework":课后作业
	// "note":笔记
	Type  string  `json:"type"`
	Msg   Msg     `json:"msg,omitempty"`
	Start float64 `json:"start,omitempty"` // 当前环节过了多长时间触发

	// 当type="turn_on_camera"时，camera_user_id才有值
	CameraUserId int64 `json:"user_id,omitempty"` // 打开摄像头的专属用户ID

	// 当type="note"|"homework"时，text才有值
	Text string `json:"text,omitempty"` // type="note | homework"时有值 "note":笔记 "homework":课后作业

	// 当type="game_input"|"game_select"时，下列所有key才有值
	Appearance         int64  `json:"appearance,omitempty"`  // 呈现风格 type="game_input" 1.我说上句你来接 2.排序题 3.课前简答题 4.课后简答题 5.喜欢的字 | type="game_select" 1.文本单选题 2.图片选择题 3.火眼金睛
	WithoutAnswerImage string `json:"background,omitempty"`  // 不带答案的互动图
	WithinAnswerImage  string `json:"background2,omitempty"` // 带答案的互动图
	Title              string `json:"title,omitempty"`       // 互动标题 eg:飞花令
	List               []Game `json:"list,omitempty"`        // 题目列表
	Task               string `json:"task,omitempty"`        // 题目描述 eg:请说出诗句的下半句
	Talks              []Talk `json:"talks,omitempty"`       // 交互过程中的会话 eg: {"say":"这题貌似讲过，我刚走神了","user_id":2}
}

// Msg - 消息内容结构体
type Msg struct {
	Mentions []int64 `json:"mentions"`  // 回复某人，支持多人
	Say      string  `json:"say"`       // 消息内容
	Tag      int64   `json:"tag"`       // 1.有感而发 (2.寒暄前 3.寒暄后 当type="TALK"时才有用，其余环节都为1)
	UserID   int64   `json:"user_id"`   // 用户ID "100":教师 | "999":助教 | 其余都为学生ID
	UserType string  `json:"user_type"` // "teacher":教师 | "assistant":助教 | "student":学生
}

// Game - 互动事件中的题目
type Game struct {
	Question  string `json:"question"`                     // 问题
	Img       string `json:"img,omitempty"`                // 图片
	Tips      string `json:"tips,omitempty"`               // 提示
	Answer    string `json:"answer"`                       // 答案
	Analyse   string `json:"analyse,omitempty"`            // 解析
	AnswerAsk string `json:"answer_requirement,omitempty"` // 答题要求
	Options   []struct {
		Content string `json:"content"`        // 选项内容
		Img     string `json:"img,omitempty"`  // 选项图片
		Tag     string `json:"tag"`            // 选项标签
		Tips    string `json:"tips,omitempty"` // 选项提示
	} `json:"options,omitempty"` // 选项列表
}

//  type Game struct {
// 	Question string   `json:"question"`          // 问题
// 	Img      string   `json:"img,omitempty"`     // 图片
// 	Tips     string   `json:"tips,omitempty"`    // 提示
// 	Answer   string   `json:"answer"`            // 答案
// 	Options  []Option `json:"options,omitempty"` // 选项
// 	Analyse  string   `json:"analyse,omitempty"` // 解析
// }

type Option struct {
	Content string `json:"content"`        // 选项内容
	Img     string `json:"img,omitempty"`  // 选项图片
	Tag     string `json:"tag"`            // 选项标签
	Tips    string `json:"tips,omitempty"` // 选项提示
}

type Talk struct {
	Say    string `json:"say"`
	UserID int64  `json:"user_id"`
}

// // Role - 角色信息
// type Role struct {
// 	RoleNo    int64  `json:"user_id"`   // 角色编号
// 	Role      string `json:"role"`      // 角色类型 主讲老师 | 班主任 | 学生
// 	RoleType  string `json:"role_type"` // "teacher":教师 | "assistant":助教 | "student":学生
// 	Character string `json:"character"` // 特征
// }
//
// // MessageList - 消息列表
// type MessageList []Message
// type Message struct {
// 	Id      int    `json:"id"`      // Id
// 	Type    string `json:"type"`    // 类型 "show-camera":打开摄像头 "chat":消息 "game":互动事件
// 	Time    int    `json:"time"`    // 当前环节过了多长时间触发，单位毫秒
// 	Uid     int    `json:"uid"`     // 用户ID
// 	GameId  int    `json:"game_id"` // 互动事件ID
// 	Content string `json:"content"` // 消息内容
// }
// type GameEvent struct {
// 	Id         int    `json:"id"`
// 	Type       int    `json:"type"`                 // 类型 1.填空 2.选择
// 	Appearance string `json:"appearance,omitempty"` // 呈现风格 type=1 f1,f2 | type=2 s1,s2,s3,s4
// 	Title      string `json:"title,omitempty"`      // 互动标题 eg:飞花令
// 	GameImage  string `json:"game_image,omitempty"` // 互动背景
// 	List       []Game `json:"list,omitempty"`       // 题目列表
// 	Task       string `json:"task,omitempty"`       // 题目描述 eg:请说出诗句的下半句
// 	Talks      []Talk `json:"talks,omitempty"`      // 交互过程中的会话 eg: {"say":"这题貌似讲过，我刚走神了","user_id":2}
// }

/* 以下为对数据源所加工的结构体 */

// Transfer 中转数据
type Transfer struct {
	CourseId          int64     // 课堂ID
	TeacherId         int64     // 老师ID
	PosterImage       string    // 视频开场背景URL，在link环节的SMALL_CHAT中
	NavigateGameStart []float64 // 环节导航栏中事件触发的时间
	FrontEvent        string    // 班主任欢迎环节
	SmallTalkEvent    string    // 寒暄环节
	SmallTalkPrecast  string    // 寒暄开场预制
	LearnEvent        string    // 课中环节预制
	// Message []Message // 弹幕消息列表
}
type Message struct {
	Id       int64   `json:"id,omitempty"`       // Id
	Type     string  `json:"type,omitempty"`     // 类型 "show-camera":打开摄像头 "chat":消息 "game":互动事件
	Time     float64 `json:"time"`               // 当前环节过了多长时间触发，单位毫秒
	UserId   int64   `json:"userId,omitempty"`   // 用户ID
	Content  string  `json:"content,omitempty"`  // 消息内容
	ToUserId int64   `json:"toUserId,omitempty"` // @用户ID
	GameId   int64   `json:"gameId,omitempty"`   // 互动事件ID
}
type PreCast struct {
	SayType string `json:"sayType,omitempty"` // 类型 "speak":发音 "text":文本
	UserId  int64  `json:"userId,omitempty"`  // 用户ID
	Say     string `json:"say,omitempty"`     // 消息内容
}
type Navigate struct {
	Name      string  `json:"name"`      // 环节名称
	StartTime float64 `json:"startTime"` // 触发时间
	Num       int64   `json:"num"`       // 在该环节中互动事件的个数
}
type Draft struct {
	Text      string `json:"text"`      // 字幕内容
	TimeBegin int64  `json:"timeBegin"` // 开始时间
	TimeEnd   int64  `json:"timeEnd"`   // 结束时间
}
type TaskTalk struct {
	Say    string `json:"say"`
	UserID int64  `json:"userId"`
}
