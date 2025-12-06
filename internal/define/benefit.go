package define

var (
	// BenefitsSource 权益来源
	BenefitsSource = struct {
		Redeem               int64 // 兑换码
		ThreePartsOfChannel  int64 // 渠道三部
		MiddleGround         int64 // 中台
		BlueVEnterpriseStore int64 // 蓝V企业店
		ProductOperation     int64 // 产品运营
		UserBuy              int64 // 用户购买
		Trial                int64 // 体验用户
		OrderCenter          int64 // 订单中心
		InternalStaff        int64 // 内部员工
	}{
		Redeem:               1, // 兑换码
		ThreePartsOfChannel:  2, // 渠道三部
		MiddleGround:         3, // 中台
		BlueVEnterpriseStore: 4, // 蓝V企业店
		ProductOperation:     5, // 产品运营
		UserBuy:              6, // 用户购买
		Trial:                7, // 体验用户
		OrderCenter:          8, // 订单中心
		InternalStaff:        9, // 内部员工
	}

	BenefitsSourceMapString = map[string]string{
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

	BenefitsSourceMapInt = map[int64]string{
		1:   "兑换码",
		2:   "渠道三部",
		3:   "中台",
		4:   "蓝V企业店",
		5:   "产品运营",
		6:   "用户购买",
		7:   "体验用户",
		8:   "订单中心",
		9:   "内部员工",
		10:  "活动",
		100: "后台手动录入",
	}

	BenefitsTypeMapInt = map[int64]string{
		1: "快解题",
		2: "关卡",
		3: "风的颜色",
		4: "学习规划",
		5: "直播课",
	}

	BenefitsStatus = struct {
		Effective int64 // 有效
		Invalid   int64 // 失效
	}{
		Effective: 1,
		Invalid:   2,
	}
)
