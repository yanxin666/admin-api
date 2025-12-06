package define

var (
	GrantStatusStr = map[int64]string{
		1: "待发货",
		2: "已发货",
		3: "已回收",
	}

	PayChannelStr = map[int64]string{
		1: "微信",
		2: "支付宝",
		3: "苹果",
		4: "余额",
	}
)
