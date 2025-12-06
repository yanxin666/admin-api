package xxl_job

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"fmt"
	"github.com/xxl-job/xxl-job-executor-go"
	"time"
)

// HelloDemo 测试跑任务demo
func HelloDemo(ctx context.Context, param *xxl.RunReq) (msg string) {
	var s string
	ctx = context.WithValue(ctx, "jobKey", param.ExecutorHandler)
	str := fmt.Sprintf("当前时间为:%s", time.Now().Format(util.StandardDatetime))
	switch expr := ctx.Value("jobKey").(type) {
	case string:
		s = expr
	default:
	}
	res := util.ConcatString(str, "正在执行的job:", s)
	return res
}
