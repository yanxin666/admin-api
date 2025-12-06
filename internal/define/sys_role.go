package define

var (
	// SysRole 系统角色
	SysRole = struct {
		Admin int64
		Demo  int64
		RD    int64
		TCH   int64
	}{
		Admin: 1, // 管理员
		Demo:  2, // 演示
		RD:    3, // 研发
		TCH:   4, // 教师
	}
)
