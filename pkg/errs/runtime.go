package errs

import (
	"github.com/spf13/cast"
	"runtime"
	"strconv"
)

// PrintCallerNameAndLine 附加信息：包名+函数名+行号
func PrintCallerNameAndLine(code BasicCode, msg string) string {
	pc, _, line, _ := runtime.Caller(2)
	codeStr := cast.ToString(int(code))
	return runtime.FuncForPC(pc).Name() + "()@" + strconv.Itoa(line) + " {" + codeStr + "：" + msg + "} >>> "
}

// PrintCallerNameAndLineByWithErr 定位记录的包名、函数名、行号
func PrintCallerNameAndLineByWithErr() string {
	pc, _, line, _ := runtime.Caller(2)
	return "追踪错误在" + runtime.FuncForPC(pc).Name() + "() line " + strconv.Itoa(line)
}
