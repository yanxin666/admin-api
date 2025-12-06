package errs

type BasicCode int

// WarnCodeListByRPC 属于Warn错误，按理全都需要改为 WarningCode。
// 但在实际场景中，端上的同学会根据某些错误枚举来进行不同的策略，故放在此列表中进行处理
var WarnCodeListByRPC = []BasicCode{
	WarningCode, // warn错误码，不会触发监控报警
}

// InfoCodeListByRPC 错误级别太低，不需要使用Warn和Err，顾为Info级别
var InfoCodeListByRPC = []BasicCode{
	InfoCode,
}

var WarnCodeListByGateway = []BasicCode{
	WarningCode, // warn错误码，不会触发监控报警
}

var InfoCodeListByGateway = []BasicCode{
	ErrUserNoLogin, // info错误码，不会触发监控报警
}

const (
	ErrCodeProgram BasicCode = 500 // 程序内部错误

	// WarningCode warning标识，使用该错误码不会触发监控报警
	WarningCode BasicCode = 510

	// InfoCode warning标识，使用该错误码不会触发监控报警
	InfoCode BasicCode = 520

	// ErrStream 900 ~ 999 流式输出错误码
	ErrStream             BasicCode = 900
	ErrCodeStreamStop     BasicCode = 901 // 流式请求输出结束
	ErrCodeFakeStreamStop BasicCode = 902 // 伪流式请求输出结束
	ErrPushMQ             BasicCode = 903 // 推送mq失败
	ErrTooFast            BasicCode = 904 // 操作太快了，等一下再试试吧

	// ServerErrorCode 1000 ~ 1099 系统错误码
	ServerErrorCode       BasicCode = 1000 // 服务繁忙，请稍后重试
	ErrCodeParamsAbnormal BasicCode = 1001 // 参数异常
	ForbiddenErrorCode    BasicCode = 1002 // 禁止操作
	NotPermMenuErrorCode  BasicCode = 1003 // 权限不足
	ErrCodeAbnormal       BasicCode = 1004 // 数据异常
	ErrSign               BasicCode = 1005 // 获取签名校验失败
	ErrSignCheck          BasicCode = 1006 // 签名校验失败
	RequestIllegal        BasicCode = 1007 // 请求非法
	OauthTokenFail        BasicCode = 1008 // 鉴权失败
	RequestConfigLack     BasicCode = 1009 // 缺少请求下游的配置项
	RequestRalFail        BasicCode = 1010 // 下游请求失败
	OperationFail         BasicCode = 1011 // 操作失败
	FileTypeFail          BasicCode = 1012 // 解析文件内容失败，请检查文件格式是否正常

	// ErrUser 1100 ~ 1199 用户错误码
	ErrUser                 BasicCode = 1100
	UserIdErrorCode         BasicCode = 1101 // 用户不存在
	CaptchaErrorCode        BasicCode = 1102 // 验证码错误
	AccountErrorCode        BasicCode = 1103 // 账号错误
	PasswordErrorCode       BasicCode = 1104 // 密码错误
	AccountDisableErrorCode BasicCode = 1105 // 账号已禁用
	AddUserErrorCode        BasicCode = 1106 // 账号已存在
	ErrUserNoLogin          BasicCode = 1107 // 账号未登录
	PassportFormatErrorCode BasicCode = 1108 // 密码格式有误
	PassportInitErrorCode   BasicCode = 1109 // 请更新初始化密码

	// ErrRole 1200 ~ 1299 角色错误码
	ErrRole                      BasicCode = 1200
	DeletePermMenuErrorCode      BasicCode = 1206 // 该权限菜单存在子级权限菜单
	ParentPermMenuErrorCode      BasicCode = 1207 // 父级菜单不能为自己
	AddRoleErrorCode             BasicCode = 1208 // 角色已存在
	DeleteRoleErrorCode          BasicCode = 1209 // 该角色存在子角色
	AddDeptErrorCode             BasicCode = 1210 // 部门已存在
	DeleteDeptErrorCode          BasicCode = 1211 // 该部门存在子部门
	AddJobErrorCode              BasicCode = 1212 // 岗位已存在
	DeleteJobErrorCode           BasicCode = 1213 // 该岗位正在使用中
	AddProfessionErrorCode       BasicCode = 1214 // 职称已存在
	DeleteProfessionErrorCode    BasicCode = 1215 // 该职称正在使用中
	DeptHasUserErrorCode         BasicCode = 1217 // 该部门正在使用中
	RoleIsUsingErrorCode         BasicCode = 1218 // 该角色正在使用中
	ParentRoleErrorCode          BasicCode = 1219 // 父级角色不能为自己
	ParentDeptErrorCode          BasicCode = 1220 // 父级部门不能为自己
	SetParentIdErrorCode         BasicCode = 1222 // 不能设置子级为自己的父级
	SetParentTypeErrorCode       BasicCode = 1223 // 权限类型不能作为父级菜单
	AddConfigErrorCode           BasicCode = 1224 // 配置已存在
	AuthErrorCode                BasicCode = 1226 // 授权已失效，请重新登录
	JobIsUsingErrorCode          BasicCode = 1228 // 该岗位正在使用中
	ProfessionIsUsingErrorCode   BasicCode = 1229 // 该职称正在使用中
	UpdateRoleUniqueKeyErrorCode BasicCode = 1231 // 角色标识已存在
	UpdateDeptUniqueKeyErrorCode BasicCode = 1232 // 部门标识已存在
	AssigningRolesErrorCode      BasicCode = 1233 // 角色不在可控范围
	DeptIdErrorCode              BasicCode = 1234 // 部门不存在
	ProfessionIdErrorCode        BasicCode = 1235 // 职称不存在
	JobIdErrorCode               BasicCode = 1236 // 岗位不存在
	RoleIdErrorCode              BasicCode = 1243 // 角色不存在
	ParentRoleIdErrorCode        BasicCode = 1237 // 父级角色不存在
	ParentDeptIdErrorCode        BasicCode = 1238 // 父级部门不存在
	ParentPermMenuIdErrorCode    BasicCode = 1239 // 父级菜单不存在
	PermMenuIdErrorCode          BasicCode = 1242 // 权限菜单不存在

	// ErrDict 1300 ~ 1399 字典错误码
	ErrDict                     BasicCode = 1300
	ParentDictionaryIdErrorCode BasicCode = 1301 // 字典集不存在
	DictionaryIdErrorCode       BasicCode = 1302 // 字典不存在
	DeleteDictionaryErrorCode   BasicCode = 1303 // 该字典集存在配置项
	AddDictionaryErrorCode      BasicCode = 1304 // 字典已存在

	// ErrMember 1400 ~ 1499 用户错误码
	ErrMember         BasicCode = 1400
	ErrMemberNoLogin  BasicCode = 1401 // 用户未登录
	ErrUserEdit       BasicCode = 1402 // 修改用户信息失败
	ErrUserHourFailed BasicCode = 1403 // 查询用户课时失败

	// ErrAppointmentConfig ErrAppointment 1500 ~ 1599 预约错误码
	ErrAppointmentConfig      BasicCode = 1500 // 配置列表查询失败
	ErrAppointmentUser        BasicCode = 1501 // 用户预约列表查询失败
	ErrUserCarrierNotSupport  BasicCode = 2018 // 不支持该类型认证
	ErrAppointmentTimeInvalid BasicCode = 1502 // 时间不合法
	ErrAppointmentAddConfig   BasicCode = 1503 // 新增预约配置失败
	ErrAppointmentEditConfig  BasicCode = 1504 // 编辑预约配置失败
	ErrAppointmentDelConfig   BasicCode = 1505 // 删除预约配置失败

	ErrCosDownloadError    BasicCode = 1601 // 下载文件出错
	ErrCreateFileFailError BasicCode = 1602 // 生成本地文件出错
	ErrWriteFileFailError  BasicCode = 1603 // 写入本地文件出错
)

// ToMsg code转换成错误信息
func ToMsg(code BasicCode) string {
	var m = map[BasicCode]string{
		ErrCodeProgram: "程序内部错误",

		// 流式输出错误码
		ErrCodeStreamStop:     "流式请求输出结束",
		ErrCodeFakeStreamStop: "伪流式请求输出结束",
		ErrPushMQ:             "导入推送mq失败",
		ErrTooFast:            "操作太快了，等一下再试试吧",

		// 系统错误码
		ServerErrorCode:         "服务繁忙，请稍后重试",
		ErrCodeParamsAbnormal:   "参数异常",
		ForbiddenErrorCode:      "禁止操作",
		NotPermMenuErrorCode:    "权限不足",
		ErrCodeAbnormal:         "数据异常",
		ErrSign:                 "获取签名校验失败",
		ErrSignCheck:            "签名校验失败",
		RequestIllegal:          "请求非法",
		OauthTokenFail:          "鉴权失败",
		RequestConfigLack:       "缺少请求下游的配置项",
		RequestRalFail:          "下游请求失败",
		OperationFail:           "操作失败",
		FileTypeFail:            "解析文件内容失败，请检查文件格式是否正常",
		CaptchaErrorCode:        "验证码错误",
		AccountErrorCode:        "账号错误",
		PasswordErrorCode:       "密码错误",
		PassportFormatErrorCode: "新密码长度在6-20个字符之间，必须同时包含字母、数字和特殊符号（如.!@#$%等）",
		PassportInitErrorCode:   "请更新初始化密码",

		DeletePermMenuErrorCode:      "该权限菜单存在子级权限菜单",
		ParentPermMenuErrorCode:      "父级菜单不能为自己",
		AddRoleErrorCode:             "角色已存在",
		DeleteRoleErrorCode:          "该角色存在子角色",
		AddDeptErrorCode:             "部门已存在",
		DeleteDeptErrorCode:          "该部门存在子部门",
		AddJobErrorCode:              "岗位已存在",
		DeleteJobErrorCode:           "该岗位正在使用中",
		AddProfessionErrorCode:       "职称已存在",
		DeleteProfessionErrorCode:    "该职称正在使用中",
		AddUserErrorCode:             "账号已存在",
		ErrUserNoLogin:               "账号未登录",
		DeptHasUserErrorCode:         "该部门正在使用中",
		RoleIsUsingErrorCode:         "该角色正在使用中",
		ParentRoleErrorCode:          "父级角色不能为自己",
		ParentDeptErrorCode:          "父级部门不能为自己",
		AccountDisableErrorCode:      "账号已禁用",
		SetParentIdErrorCode:         "不能设置子级为自己的父级",
		SetParentTypeErrorCode:       "权限类型不能作为父级菜单",
		AddConfigErrorCode:           "配置已存在",
		AddDictionaryErrorCode:       "字典已存在",
		AuthErrorCode:                "授权已失效，请重新登录",
		DeleteDictionaryErrorCode:    "该字典集存在配置项",
		JobIsUsingErrorCode:          "该岗位正在使用中",
		ProfessionIsUsingErrorCode:   "该职称正在使用中",
		UpdateRoleUniqueKeyErrorCode: "角色标识已存在",
		UpdateDeptUniqueKeyErrorCode: "部门标识已存在",
		AssigningRolesErrorCode:      "角色不在可控范围",
		DeptIdErrorCode:              "部门不存在",
		ProfessionIdErrorCode:        "职称不存在",
		JobIdErrorCode:               "岗位不存在",
		ParentRoleIdErrorCode:        "父级角色不存在",
		ParentDeptIdErrorCode:        "父级部门不存在",
		ParentPermMenuIdErrorCode:    "父级菜单不存在",
		ParentDictionaryIdErrorCode:  "字典集不存在",
		DictionaryIdErrorCode:        "字典不存在",
		PermMenuIdErrorCode:          "权限菜单不存在",
		RoleIdErrorCode:              "角色不存在",
		UserIdErrorCode:              "用户不存在",

		// 用户错误
		ErrMemberNoLogin:  "用户未登录",
		ErrUserEdit:       "修改用户信息失败",
		ErrUserHourFailed: "查询用户课时失败",

		// 预约错误码
		ErrAppointmentConfig:      "查询配置列表失败",
		ErrAppointmentUser:        "查询用户预约列表失败",
		ErrAppointmentTimeInvalid: "预约配置时间不合法", // 时间不合法
		ErrAppointmentAddConfig:   "新增预约配置失败",
		ErrAppointmentEditConfig:  "更新预约配置失败",
		ErrAppointmentDelConfig:   "删除预约配置失败",
	}
	if msg, ok := m[code]; ok {
		return msg
	}
	return "服务繁忙，请稍后重试"
}
