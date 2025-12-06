package define

var (
	// OpinionStatus 意见反馈状态
	OpinionStatus = struct {
		OpinionStatusNew int64
		opinionStatusOld int64
	}{
		OpinionStatusNew: 0, // 新意见
		opinionStatusOld: 1, // 已回复
	}

	// IsInfoFinishedStatus 用户完善用户状态
	IsInfoFinishedStatus = struct {
		Unfinished int64
		Finish     int64
	}{
		Unfinished: 1, // 未完善
		Finish:     2, // 已完善
	}

	// IsAssessment 用户测评状态
	IsAssessment = struct {
		Unfinished int64
		Finish     int64
	}{
		Unfinished: 1, // 未评测
		Finish:     2, // 已评测
	}

	// IsGuide 用户引导状态
	IsGuide = struct {
		No  int64
		Yes int64
	}{
		No:  0, // 未引导
		Yes: 1, // 已引导
	}

	// IsGuideTips 引导卡片下方Tips是否展示过
	IsGuideTips = struct {
		No  int64
		Yes int64
	}{
		No:  0, // 未展示
		Yes: 1, // 已展示
	}

	// GenderType 性别
	GenderType = struct {
		Male    int64
		Female  int64
		Unknown int64
	}{
		Unknown: 0, // 未知
		Male:    1, // 男
		Female:  2, // 女
	}

	// GenderTypePersonMap 性别人称映射
	GenderTypePersonMap = map[int64]string{
		GenderType.Unknown: "ta",
		GenderType.Male:    "他",
		GenderType.Female:  "她",
	}

	// UserRole 用户角色
	UserRole = struct {
		Child  int64
		Parent int64
	}{
		Child:  1, // 孩子
		Parent: 2, // 家长
	}
)

// GetUserGenderInfo 获取用户性别信息
func GetUserGenderInfo(gender int64) string {
	if gender == GenderType.Male {
		return "男"
	} else {
		return "女"
	}
}

// GetUserTaInfo 获取用户称呼信息
func GetUserTaInfo(gender int64) string {
	if gender == GenderType.Male {
		return "他"
	} else {
		return "她"
	}
}

// GetUserRole 获取用户角色
func GetUserRole(role int64) string {
	if role == UserRole.Parent {
		return "家长"
	} else {
		return "孩子"
	}
}
