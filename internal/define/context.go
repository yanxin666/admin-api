package define

// ContextKey 上下文传导value值对应的key
type ContextKey string

const (
	CtxUserID        string = "user_id"        // 用户ID
	CtxUser          string = "user"           // 用户信息
	CtxToken         string = "token"          // 用户token
	CtxCommonParams  string = "common_params"  // 通用参数
	CtxRequestParams string = "request_params" // 请求数据
	CtxRequestLogs   string = "request_logs"   // 请求日志
	CtxTrace         string = "trace"          // 请求日志
)
