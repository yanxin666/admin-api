package define

// 用户相关
const (
	RedisKeyUserInfoKey       = "api:user:userInfo:%d" // 用户缓存
	RedisKeyUserInfoKeyExpire = 86400                  // 用户缓存失效时间-1天 单位(s)
)

// LLM大模型相关
const (
	// RedisKeyResourceList 大模型资源列表缓存key
	RedisKeyResourceList = "api:json_resource_list:%s"
	// RedisKeyResourceListExpire 大模型资源列表缓存有效期-5分钟，单位(s)
	RedisKeyResourceListExpire = 300

	// RedisKeyResourceChoice 大模型资源列表轮询缓存key
	RedisKeyResourceChoice = "api:str_resource_choice:%s"
	// RedisKeyResourceChoiceExpire 大模型资源列表轮询缓存有效期-24小时，单位(s)
	RedisKeyResourceChoiceExpire = 86400
)

// 评测相关
const (
	// RedisKeyAssessmentRecordCreate 生成评测记录分布式锁，占位符：用户ID
	RedisKeyAssessmentRecordCreate = "api:assessment:record:create:%d"
	// RedisKeyAssessmentRecordCreateExpire 生成评测记录分布式锁有效期，10秒
	RedisKeyAssessmentRecordCreateExpire = 10

	// RedisKeyAssessmentQuestion 评测题目缓存key，占位符：题目ID
	RedisKeyAssessmentQuestion = "api:assessment:question:%d"
	// RedisKeyAssessmentQuestionExpire 评测题目缓存有效期，24h
	RedisKeyAssessmentQuestionExpire = 86400
)

// 定义引导规则相关Redis Key常量
var (
	// GuideRuleRedisKey 引导规则Redis Key常量
	GuideRuleRedisKey = struct {
		Frequency      string // 超频规则
		FrequencyMatch string // 超频规则模糊匹配
		Wrong          string // 错误次数
		WrongMatch     string // 错误次数模糊匹配
		RecordAbility  string // 一次评测中一个能力的总（题）数
	}{
		Frequency:      "api:answer:total:%d:%d",   // 超频规则，api:answer:total:${assessment_record_id}:${question_id}
		FrequencyMatch: "api:answer:total:%d:*",    // 超频规则模糊匹配，api:answer:total:${assessment_record_id}:*
		Wrong:          "api:answer:wrong:%d:%d",   // 错误次数，api:answer:wrong:${assessment_record_id}:${question_id}
		WrongMatch:     "api:answer:wrong:%d:*",    // 错误次数，api:answer:wrong:${assessment_record_id}:*
		RecordAbility:  "api:answer:ability:%d:%s", // 格式：api:answer:ability:${assessment_record_id}:${ability_code}
	}
)

const (
	// RedisKeyReportCreate 生成报告幂等
	RedisKeyReportCreate       = "api:report:create:%v"
	RedisKeyReportCreateExpire = 30

	// RedisKeyReportBaseAbility 基础能力
	RedisKeyReportBaseAbility       = "api:report:base:ability"
	RedisKeyReportBaseAbilityExpire = 604800 // 7天

	// RedisKeyReportSortBaseAbility 基础能力(排序后)(T1、T2、T3)
	RedisKeyReportSortBaseAbility       = "api:report:sort:base:ability"
	RedisKeyReportSortBaseAbilityExpire = 604800 // 7天

	// RedisKeyReportSortDiffBaseAbility 基础能力(排序后)(难度)
	RedisKeyReportSortDiffBaseAbility       = "api:report:sort:diff:base:ability"
	RedisKeyReportSortDiffBaseAbilityExpire = 604800 // 7天
)

// 定义预约课程规则相关Redis Key常量
var (
	// AppointmentRedisKey 预约课程Redis Key常量
	AppointmentRedisKey = struct {
		AppointSetNx  string // 分布式锁
		Booked        string // 已预约课程ID集合
		Stock         string // 课程库存
		CurrentUsable string // 用户当前可用的预约缓存
	}{
		AppointSetNx:  "api:appointment:setnx:%d_%d", // 分布式锁：api:appointment:setnx:${user_id}_${appointment_config_id}
		Booked:        "api:appointment:booked:%d",   // 已预约课程ID集合，集合中存储的是用户ID：api:appointment:booked:${appointment_config_id}
		Stock:         "api:appointment:stock:%d",    // 课程库存：api:appointment:stock:${appointment_config_id}
		CurrentUsable: "api:appointment:usable:%d",   // 用户当前可用的预约缓存：api:appointment:${user_id}
	}
)

// 定义训练规则相关Redis Key常量
var (
	// TrainRedisKey 预约课程Redis Key常量
	TrainRedisKey = struct {
		PreGenSetNx string // 用户训练所用分布式锁
	}{
		PreGenSetNx: "api:train:pre_gen_setnx:%d", // 用户训练所用分布式锁：api:train:pre_gen_setnx:${user_id}
	}
)
