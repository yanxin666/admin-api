package define

var (
	// AppointmentUserStatus 用户预约状态
	AppointmentUserStatus = struct {
		Available int64
		Used      int64
		Expired   int64
		Cancel    int64
	}{
		Available: 1, // 可用
		Used:      2, // 已使用
		Expired:   3, // 已过期
		Cancel:    4, // 已取消
	}

	// UserAppointSearchStatus 课程管理列表搜素状态
	UserAppointSearchStatus = struct {
		ImmediateClass int64 // 立即上课
		UpcomingClass  int64 // 即将上课
		NoBooking      int64 // 暂无约课
		HistoryReview  int64 // 历史回顾
		UnKnown        int64
	}{

		ImmediateClass: 1, // 立即上课
		UpcomingClass:  2, // 即将上课
		NoBooking:      3, // 暂无约课
		HistoryReview:  4, // 历史回顾
		UnKnown:        0, // 缺省值
	}

	// AppointmentConfigStatus 预约配置状态
	AppointmentConfigStatus = struct {
		Disabled int64
		Enabled  int64
		Delisted int64
	}{
		Disabled: 1, // 未启用
		Enabled:  2, // 启用
		Delisted: 3, // 已下架
	}

	// AppointmentUserClassRecordStatus 用户课时记录状态
	AppointmentUserClassRecordStatus = struct {
		Present     int64
		Buy         int64
		Use         int64
		Return      int64
		Expired     int64
		TransferIn  int64
		TransferOut int64
	}{
		Present:     1, // 赠送
		Buy:         2, // 购买
		Use:         3, // 使用
		Return:      4, // 退回
		Expired:     5, // 过期
		TransferIn:  6, // 转入
		TransferOut: 7, // 转出

	}
)

const (
	AdvanceTimeMinute = 15 // 具体开课提前时间15分钟
	DefaultPageSize   = 10 // 默认一页10条
)

var (
	// AppointUserBindStatus 用户预约表教师绑定状态
	AppointUserBindStatus = struct {
		Unbind int64
		Bind   int64
	}{
		Unbind: 2, // 未绑定
		Bind:   1, // 已绑定
	}

	// AppointUserStatus 用户预约表状态
	AppointUserStatus = struct {
		Available int64
		Used      int64
		Expired   int64
		Cancelled int64
	}{
		Available: 1, // 可用
		Used:      2, // 已使用
		Expired:   3, // 过期
		Cancelled: 4, // 已取消
	}
)
