package define

// 豆神王者Club课程相关常量定义

// 系列类型
const (
	KCSeriesTypeChaoLianZuoWen = 1 // 超练作文
	KCSeriesTypeChaoLianYueDu  = 2 // 超练阅读
)

var KCSeriesTypeNameMap = map[int]string{
	KCSeriesTypeChaoLianZuoWen: "超练作文",
	KCSeriesTypeChaoLianYueDu:  "超练阅读",
}

// 课程类型
const (
	KCCourseTypeZuoWenManLian  = 1 // 超练作文慢练
	KCCourseTypeZuoWenKuaiLian = 2 // 超练作文快练
	KCCourseTypeYueDu          = 3 // 超练阅读
)

var KCCourseTypeNameMap = map[int]string{
	KCCourseTypeZuoWenManLian:  "超练作文慢练",
	KCCourseTypeZuoWenKuaiLian: "超练作文快练",
	KCCourseTypeYueDu:          "超练阅读",
}

// 任务类型
const (
	KCTaskTypeNormalGroup  = 1   // 普通线索组合
	KCTaskTypeIdentityCard = 2   // 身份卡
	KCTaskTypeArticleSum   = 3   // 文章总结
	KCTaskTypeTimeCorridor = 4   // 时空长廊（组合）
	KCTaskTypeMainLine     = 200 // 主线任务
	KCTaskTypeEasyCopy     = 201 // 简单副本任务
	KCTaskTypeHardCopy     = 202 // 困难副本
)

var KCTaskTypeNameMap = map[int]string{
	KCTaskTypeNormalGroup:  "普通线索组合",
	KCTaskTypeIdentityCard: "身份卡",
	KCTaskTypeArticleSum:   "文章总结",
	KCTaskTypeTimeCorridor: "时空长廊（组合）",
	KCTaskTypeMainLine:     "主线任务",
	KCTaskTypeEasyCopy:     "简单副本任务",
	KCTaskTypeHardCopy:     "困难副本",
}

// 子任务类型
const (
	KCSubTaskTypeNormal   = 1 // 普通任务
	KCSubTaskTypeIdentity = 2 // 身份卡
	KCSubTaskTypeChat     = 3 // 聊天
	KCSubTaskTypeSubtitle = 4 // 结尾字幕
	KCSubTaskTypeSignWall = 5 // 签名墙
	KCSubTaskTypeTimeShow = 6 // 时光长廊-展示
)

var KCSubTaskTypeNameMap = map[int]string{
	KCSubTaskTypeNormal:   "普通任务",
	KCSubTaskTypeIdentity: "身份卡",
	KCSubTaskTypeChat:     "聊天",
	KCSubTaskTypeSubtitle: "结尾字幕",
	KCSubTaskTypeSignWall: "签名墙",
	KCSubTaskTypeTimeShow: "时光长廊-展示",
}

// KCSeriesTypeList 系列类型列表
var KCSeriesTypeList = []map[string]interface{}{
	{
		"series_type":      KCSeriesTypeChaoLianZuoWen,
		"series_type_name": KCSeriesTypeNameMap[KCSeriesTypeChaoLianZuoWen],
	},
	{
		"series_type":      KCSeriesTypeChaoLianYueDu,
		"series_type_name": KCSeriesTypeNameMap[KCSeriesTypeChaoLianYueDu],
	},
}

// KCCourseTypeList 课程类型列表
var KCCourseTypeList = []map[string]interface{}{
	{
		"course_type":      KCCourseTypeZuoWenManLian,
		"course_type_name": KCCourseTypeNameMap[KCCourseTypeZuoWenManLian],
	},
	{
		"course_type":      KCCourseTypeZuoWenKuaiLian,
		"course_type_name": KCCourseTypeNameMap[KCCourseTypeZuoWenKuaiLian],
	},
	{
		"course_type":      KCCourseTypeYueDu,
		"course_type_name": KCCourseTypeNameMap[KCCourseTypeYueDu],
	},
}

// KCTaskTypeList 任务类型列表
var KCTaskTypeList = []map[string]interface{}{
	{
		"task_type":      KCTaskTypeNormalGroup,
		"task_type_name": KCTaskTypeNameMap[KCTaskTypeNormalGroup],
	},
	{
		"task_type":      KCTaskTypeIdentityCard,
		"task_type_name": KCTaskTypeNameMap[KCTaskTypeIdentityCard],
	},
	{
		"task_type":      KCTaskTypeArticleSum,
		"task_type_name": KCTaskTypeNameMap[KCTaskTypeArticleSum],
	},
	{
		"task_type":      KCTaskTypeTimeCorridor,
		"task_type_name": KCTaskTypeNameMap[KCTaskTypeTimeCorridor],
	},
	{
		"task_type":      KCTaskTypeMainLine,
		"task_type_name": KCTaskTypeNameMap[KCTaskTypeMainLine],
	},
	{
		"task_type":      KCTaskTypeEasyCopy,
		"task_type_name": KCTaskTypeNameMap[KCTaskTypeEasyCopy],
	},
	{
		"task_type":      KCTaskTypeHardCopy,
		"task_type_name": KCTaskTypeNameMap[KCTaskTypeHardCopy],
	},
}

// KCSubTaskTypeList 子任务类型列表
var KCSubTaskTypeList = []map[string]interface{}{
	{
		"sub_task_type":      KCSubTaskTypeNormal,
		"sub_task_type_name": KCSubTaskTypeNameMap[KCSubTaskTypeNormal],
	},
	{
		"sub_task_type":      KCSubTaskTypeIdentity,
		"sub_task_type_name": KCSubTaskTypeNameMap[KCSubTaskTypeIdentity],
	},
	{
		"sub_task_type":      KCSubTaskTypeChat,
		"sub_task_type_name": KCSubTaskTypeNameMap[KCSubTaskTypeChat],
	},
	{
		"sub_task_type":      KCSubTaskTypeSubtitle,
		"sub_task_type_name": KCSubTaskTypeNameMap[KCSubTaskTypeSubtitle],
	},
	{
		"sub_task_type":      KCSubTaskTypeSignWall,
		"sub_task_type_name": KCSubTaskTypeNameMap[KCSubTaskTypeSignWall],
	},
	{
		"sub_task_type":      KCSubTaskTypeTimeShow,
		"sub_task_type_name": KCSubTaskTypeNameMap[KCSubTaskTypeTimeShow],
	},
}

// KCTeacherIdList 教师ID列表 TAG: 临时，需要从数据库查询
var KCTeacherIdList = []map[string]interface{}{
	{
		"id":   1,
		"name": "CL001窦昕",
	},
	{
		"id":   2,
		"name": "CL003",
	},
}

// KCCourseUseCases 王者Club课程使用场景 1.练习 2.直播互动
var KCCourseUseCases = struct {
	Practice        int64
	LiveInteraction int64
}{
	Practice:        1, // 练习
	LiveInteraction: 2, // 直播互动
}

// KCScheduleCourseStatus 课程排课状态 (状态:1.待开放2.开放3.进行中4.结束)
var KCScheduleCourseStatus = struct {
	Pending  int64
	Open     int64
	Ongoing  int64
	Finished int64
}{
	Pending:  1, // 待开放
	Open:     2, // 开放
	Ongoing:  3, // 进行中
	Finished: 4, // 结束
}
