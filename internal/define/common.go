package define

var (
	// BaseUserStatus 全局所使用状态
	BaseUserStatus = struct {
		Normal int64
		Freeze int64
		LogOff int64
	}{
		Normal: 0, // 正常
		Freeze: 1, // 冻结
		LogOff: 2, // 删除
	}

	// UserStatus 全局所使用状态
	UserStatus = struct {
		Normal int64
		LogOff int64
	}{
		Normal: 1, // 正常
		LogOff: 2, // 删除
	}

	// MediumType 媒体类型
	MediumType = struct {
		Text    int64
		Picture int64
		Audio   int64
		Video   int64
	}{
		Text:    0, // 文本
		Picture: 1, // 图片
		Audio:   2, // 音频
		Video:   3, // 视频
	}

	// UserProduct 用户产品线
	UserProduct = struct {
		CiYuan     int64
		DouZiRobot int64
	}{
		CiYuan:     0, // 辞源
		DouZiRobot: 4, // 豆子机器人
	}
)

var (
	BaseUserStatusMap = map[int64]string{
		BaseUserStatus.Normal: "正常",
		BaseUserStatus.Freeze: "冻结",
		BaseUserStatus.LogOff: "删除",
	}

	UserStatusMap = map[int64]string{
		UserStatus.Normal: "正常",
		UserStatus.LogOff: "删除",
	}

	UserProductMap = map[int64]string{
		0: "辞源",
		1: "AI搜索",
		2: "白帝辞",
		3: "智普乌托邦",
	}

	UserSourceMap = map[int64]string{
		100:  "提审标记",
		0:    "APP",
		1:    "豆伴匠",
		2:    "风的颜色",
		1001: "听力熊",
		1002: "微软",
	}
)
