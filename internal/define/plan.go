package define

var (
	// PlanOutlineType 规划大纲内容类型
	PlanOutlineType = struct {
		Text  int64
		Lists int64
		Audio int64
		Video int64
	}{
		Text:  0, // 文本
		Lists: 1, // 列表
		Audio: 2, // 音频
		Video: 3, // 视频
	}
)
