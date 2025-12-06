package errs

import (
	"github.com/pkg/errors"
)

// WithErr 仅仅在日志里记录抛出错误的位置，并无它用。一般来做底层抛出err，业务需要定位位置时使用，不可再返回API层调用。
func WithErr(err error) ErrInterface {
	if err == nil {
		err = errors.New(ToMsg(ErrCodeProgram))
	}
	var (
		target  ErrInterface
		code    BasicCode
		message string
	)
	if errors.As(err, &target) {
		code = target.Code()
		message = target.Msg()
	} else {
		code = ErrCodeProgram
		message = ToMsg(ErrCodeProgram)
	}
	return &BasicErr{
		error:    err,
		code:     code,
		msg:      message,
		stackMsg: errors.WithMessage(err, PrintCallerNameAndLine(code, message)).Error(),
	}
}

// NewCode 在日志里添加 BasicCode.ToMsg 中的错误提示。若需要在日志里显示新提示，请使用 NewMsg 方法处理
func NewCode(code BasicCode) ErrInterface {
	msg := ToMsg(code)
	return &BasicErr{
		error:    errors.New(msg),
		code:     code,
		msg:      msg,
		stackMsg: errors.New(PrintCallerNameAndLine(code, msg)).Error(),
	}
}

// NewMsg 在日志里显示新提示，方便错误排查。若不需要新提示，请使用 NewCode 方法处理
// 注意：此处的msg是为了给RD方便在日志里直观的查看问题使用，而非在接口层返回！！！如需查看接口返回值，请移步下方 GetErrCodeAndMessage 方法
func NewMsg(code BasicCode, msg string) ErrInterface {
	return &BasicErr{
		error:    errors.New(msg),
		code:     code,
		msg:      msg,
		stackMsg: errors.New(PrintCallerNameAndLine(code, msg)).Error(),
	}
}

// WithCode 在已有错误基础上，在日志里添加 BasicCode.ToMsg 中的错误提示。若需要在日志里显示新提示，请使用 WithMsg 方法处理
func WithCode(err error, code BasicCode) ErrInterface {
	msg := ToMsg(code)
	if err == nil {
		err = errors.New(msg)
	}
	return &BasicErr{
		error:    err,
		code:     code,
		msg:      msg,
		stackMsg: errors.WithMessage(err, PrintCallerNameAndLine(code, msg)).Error(),
	}
}

// WithMsg 在已有错误基础上，在日志中添加新提示信息，若不需要日志里有新提示，请使用 WithCode 方法处理
// 注意：此处的msg是为了给RD方便在日志里直观的查看问题使用，而非在接口层返回！！！如需查看接口返回值，请移步下方 GetErrCodeAndMessage 方法
func WithMsg(err error, code BasicCode, msg string) ErrInterface {
	if err == nil {
		err = errors.New(msg)
	}
	return &BasicErr{
		error:    err,
		code:     code,
		msg:      msg,
		stackMsg: errors.WithMessage(err, PrintCallerNameAndLine(code, msg)).Error(),
	}
}

// GetErrCodeAndMessage 在接口响应时，根据上方调用的error来获取对应的错误码和错误信息进行返回
// 接口返回的是code对应map的值，如其值不存在，就取默认：服务繁忙
func GetErrCodeAndMessage(err error) (code BasicCode, message string) {
	// 不同类型处理
	var (
		e       ErrInterface
		codeMsg string // code映射的错误信息
		showMsg string // 自定义的错误信息
	)

	// 默认错误
	code = ErrCodeProgram
	baseMessage := ToMsg(code)
	if !errors.As(err, &e) {
		return
	}

	// 获取业务中的错误及映射的错误信息
	code = e.Code()
	if ToMsg(code) != "" {
		codeMsg = ToMsg(code)
	}

	// 获取showMsg的错误信息
	if e.GetShow() {
		showMsg = e.Msg()
	}

	// 优先展示showMsg
	if showMsg != "" {
		message = showMsg
	} else {
		message = codeMsg
	}

	// 兜底错误
	if message == "" {
		message = baseMessage
	}
	return
}

// IsBasicError 比较两个code是否一致，一般用于前后调用的关系依赖
func IsBasicError(err, target error) bool {
	var e1 *BasicErr
	ok1 := errors.As(err, &e1)
	var e2 *BasicErr
	ok2 := errors.As(target, &e2)
	if !ok1 || !ok2 {
		return false
	}
	return e1.Code() == e2.Code()
}
