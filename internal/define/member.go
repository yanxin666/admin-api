package define

// 会员相关枚举请用如下方式来定义
var (
	VipType = struct {
		DouShenVip        int64 // 豆神AI会员（基础版）
		DouShenVipPlus    int64 // 豆神AI会员（典藏版）
		DouShenSuperVip   int64 // 豆神AI会员（超级版）
		DouShenVipSpecial int64 // 豆神AI会员（特别版）
	}{
		DouShenVip:        1,
		DouShenVipPlus:    2,
		DouShenSuperVip:   3,
		DouShenVipSpecial: 4,
	}
	VipTypeMap = map[int64]string{
		VipType.DouShenVip:        "豆神AI会员(基础版)",
		VipType.DouShenVipPlus:    "豆神AI会员(典藏版)",
		VipType.DouShenSuperVip:   "豆神AI超级会员",
		VipType.DouShenVipSpecial: "豆神AI会员(特别版)",
	}

	// MemberOpinionStatus 意见反馈状态
	MemberOpinionStatus = struct {
		New int64
		Old int64
	}{
		New: 0, // 新意见
		Old: 1, // 已回复
	}

	// MemberEvalStatus 用户测评状态
	MemberEvalStatus = struct {
		Unfinished int64
		Finish     int64
	}{
		Unfinished: 0, // 未完成
		Finish:     1, // 已完成
	}
)
