package errs

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"runtime"
	"testing"
)

func testNewCode() error {
	return NewCode(ErrCodeParamsAbnormal)
}
func testNewMsg() error {
	return NewMsg(ErrCodeParamsAbnormal, "不在错误集合里定义的新提示")
}
func TestNewErr(t *testing.T) {

	err := testNewCode()
	fmt.Printf("NewCode打印错误：%v\n", err)

	err1 := testNewMsg()
	fmt.Printf("NewMsg打印错误：%v\n", err1)
}

func testWithErr() error {
	err := testWithMsg()
	return WithErr(err)
}
func testWithMsg() error {
	err := testWithCode()
	return WithMsg(err, ErrAppointmentConfig, "不在错误集合里定义的新提示")
}
func testWithCode() error {
	err := testNeedErr()
	return WithCode(err, ErrRole)
}
func testNeedErr() error {
	err := sqlx.ErrNotFound
	// 若想要在日志里体现 具体带行数的报错，请用 WithErr 方法
	// return WithErr(err)
	return err
}
func TestWithErr(t *testing.T) {
	err := testWithErr()
	fmt.Printf("WithErr打印错误：%v\n", err)

	code, message := GetErrCodeAndMessage(err)
	fmt.Printf("接口返回数据：%v  %v\n", code, message)

	// 最终返回的日志
	_, file, line, _ := runtime.Caller(1)
	logc.Errorf(context.Background(), "接口报错：%s:%d，错误信息：%s", file, line, err)
}
