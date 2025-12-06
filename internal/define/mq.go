package define

var (
	// ProducerTopics 需要发送的MQ
	ProducerTopics = map[string]struct {
		ExcelImport          string // excel导入数据生成任务
		ExcelExport          string // excel导出数据生成任务
		UserBenefitPreImport string // 用户权益预导入处理
		UserBalanceImport    string // 用户金额导入处理
	}{
		"dev": {
			ExcelImport:          "muse_admin_import_dev",
			ExcelExport:          "zx_admin_export",
			UserBenefitPreImport: "user_benefit_pre_import_dev",
			UserBalanceImport:    "muse_ability_balance_import_dev",
		},
		"test": {
			ExcelImport:          "muse_admin_import_test",
			ExcelExport:          "zx_admin_export",
			UserBenefitPreImport: "user_benefit_pre_import_test",
			UserBalanceImport:    "muse_ability_balance_import_test",
		},
		"pre": {
			ExcelImport:          "muse_admin_import_dev",
			ExcelExport:          "zx_admin_export",
			UserBenefitPreImport: "user_benefit_pre_import_pre",
			UserBalanceImport:    "muse_ability_balance_import_pre",
		},
		"pro": {
			ExcelImport:          "muse_admin_import",
			ExcelExport:          "zx_admin_export",
			UserBenefitPreImport: "user_benefit_pre_import",
			UserBalanceImport:    "muse_ability_balance_import",
		},
	}

	// ConsumerTopic 需要消费的MQ
	ConsumerTopic = map[string]struct {
		ExcelImport       string // excel导入数据生成任务
		ExcelExport       string // excel导出数据生成任务
		ScheduleKnowledge string // 规划题库数据
	}{
		"dev": {
			ExcelImport:       "muse_admin_import_dev",
			ExcelExport:       "zx_admin_export",
			ScheduleKnowledge: "jx_admin_lecture_dev",
		},
		"test": {
			ExcelImport:       "muse_admin_import_test",
			ExcelExport:       "zx_admin_export",
			ScheduleKnowledge: "jx_admin_lecture_test",
		},
		"pre": {
			ExcelImport:       "muse_admin_import_dev",
			ExcelExport:       "zx_admin_export",
			ScheduleKnowledge: "jx_admin_lecture_dev",
		},
		"pro": {
			ExcelImport:       "muse_admin_import",
			ExcelExport:       "zx_admin_export",
			ScheduleKnowledge: "jx_admin_lecture",
		},
	}

	// OnlyTestTopicGroup 仅测试环境组监听的Topic
	OnlyTestTopicGroup = []string{
		HubConsumerToTest,
	}

	// OnlyProTopicGroup 仅生产环境组监听的Topic
	OnlyProTopicGroup = []string{
		HubConsumerToPro,
	}
)

// 数据中心
const (
	GroupHubConsumer  = "group_hub_consumer"   // Group 接收由【线上环境】向【测试环境】发送数据进行消费的消费者组
	HubConsumerToTest = "hub_consumer_to_test" // 线上推送测试，直接更新测试数据
	HubConsumerToPro  = "hub_consumer_to_pro"  // 审核通过后，需要更新的线上数据
)
