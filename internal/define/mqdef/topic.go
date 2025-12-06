package mqdef

import (
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/mq/types"
)

const (
	/*
		这里是基础的Topic定义，如果需要在这里添加新的Topic，请按照以下格式添加
	*/
	TopicTest types.Topic = "muse_ability_simple" // 测试Topic

	TopicUserBenefitPreImport types.Topic = "user_benefit_pre_import"     // 用户权益预导入处理
	TopicUserBalanceImport    types.Topic = "muse_ability_balance_import" // 用户金额导入处理

	TopicScheduleKnowledge types.Topic = "jx_admin_lecture"         // 规划题库数据
	TopicHubConsumerToTest types.Topic = "hub_consumer_online2test" // 线上推送测试，直接更新测试数据
	TopicHubConsumerToPre  types.Topic = "hub_consumer_online2pre"  // 线上推送预发，直接更新预发数据
	TopicHubConsumerToPro  types.Topic = "hub_consumer_to_pro"      // 审核通过后，需要更新的线上数据
)
