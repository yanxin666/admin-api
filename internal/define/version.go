package define

var (
	// VersionConfigStatus 版本配置状态
	VersionConfigStatus = struct {
		// 未发布
		Unpublished int64
		// 已发布
		Published int64
		// 已下线
		Offline int64
	}{
		Unpublished: 0,
		Published:   1,
		Offline:     2,
	}
)
