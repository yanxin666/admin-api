package errs

// ErrInterface Err接口和实体
type ErrInterface interface {
	Raw() error
	Code() BasicCode
	Msg() string
	Error() string
	GetShow() bool
	ShowMsg() ErrInterface
	Warning() ErrInterface // 该方法用处：为防止基础错误结束流程报错而报警。一般用在非系统级别错误上，比如：参数异常、鉴权失败、验证码效验失败等业务错误，不推送报警群。
}

type BasicErr struct {
	error
	code     BasicCode
	msg      string
	stackMsg string
	show     bool
}

// Raw 获取error
func (e *BasicErr) Raw() error {
	return e.error
}

// Code 获取最新的错误码
func (e *BasicErr) Code() BasicCode {
	return e.code
}

// Msg 获取最新的错误信息
func (e *BasicErr) Msg() string {
	return e.msg
}

// StackMsg 获取全链路错误信息
func (e *BasicErr) Error() string {
	return e.stackMsg
}

// GetShow 获取最新的错误显示标识
func (e *BasicErr) GetShow() bool {
	return e.show
}

// ShowMsg 修改 show 为 true，使msg在API可见
func (e *BasicErr) ShowMsg() ErrInterface {
	e.show = true
	return e
}

// Warning 标识此错误为Warning错误，不进行报警处理
func (e *BasicErr) Warning() ErrInterface {
	e.stackMsg = PrintCallerNameAndLine(WarningCode, e.msg)
	// 若无showMsg标识，优先展示code映射内容
	if !e.show {
		e.msg = ToMsg(e.Code())
	}
	e.code = WarningCode
	return e
}
