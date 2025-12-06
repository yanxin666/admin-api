package _import

const (
	RedisKeyImportCount = "muse-admin:import:count:%d" // 当前任务执行数量
)

var (
	StatusType = struct {
		NotStart int64
		Success  int64
		Going    int64
		Fail     int64
	}{
		NotStart: 0, // 未开始
		Going:    1, // 进行中
		Success:  2, // 成功
		Fail:     3, // 失败
	}
)

var (
	bizTypeMap = map[string]string{
		"评测":   "1",
		"爬天梯":  "2",
		"原理讲解": "3",
		"精析拔高": "4",
		"冲刺高考": "5",
	}

	usageTypeMap = map[string]string{
		"例题":  "1",
		"练习题": "2",
		"候补题": "3",
	}

	questionTypeMap = map[string]string{
		"单选": "1",
		"多选": "2",
		"填空": "3",
		"判断": "4",
		"简答": "5",
		"阅读": "6",
		"作文": "7",
	}

	orderStatusMap = map[string]string{
		"已完成":  "1",
		"部分退款": "3",
		"全额退款": "2",
	}

	sourceMap = map[string]string{
		"兑换码":   "1",
		"渠道三部":  "2",
		"中台":    "3",
		"蓝V企业店": "4",
		"产品运营":  "5",
		"用户购买":  "6",
		"体验用户":  "7",
		"订单中心":  "8",
		"内部员工":  "9",
	}
)
