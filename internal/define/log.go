package define

var (
	// LogFieldsType 日志字段类型
	LogFieldsType = struct {
		Path   string
		UserID string
	}{
		Path:   "path",   // 请求路径
		UserID: "userId", // 用户ID
	}

	// SysLogType 系统操作日志类型
	SysLogType = struct {
		Login    int64
		Operator int64
	}{
		Login:    1, // 登录type
		Operator: 2, // 操作type
	}
)
