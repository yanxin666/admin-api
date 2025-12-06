package tools

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/logz"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"github.com/zeromicro/go-zero/core/logc"
	"muse-admin/pkg/errs"
	"runtime"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// rpcErrCode 需要处理的RPC错误码
var rpcErrCode = []codes.Code{
	codes.Unavailable, // 未连接到RPC / 连接被拒绝
	codes.Internal,    // 内部错误
}

// RpcErrChange RPC返回的error转换
func RpcErrChange(ctx context.Context, err error) (code errs.BasicCode, message string) {
	code = errs.ErrCodeProgram
	message = errs.ToMsg(code)

	if errStatus, ok := status.FromError(err); ok {
		rpcCode := errStatus.Code()
		code = errs.BasicCode(rpcCode)
		// 需要收纳的RPC错误，比如下游连接异常或超时等系统错误时，需返回默认兜底错误
		if !util.IsExist(rpcCode, rpcErrCode) {
			message = errStatus.Message()
		}
		// 打印不同错误级别的日志
		if util.IsExist(code, errs.InfoCodeListByRPC) {
			logz.Infof(ctx, "rpc error code: %d，error message: %s，error: %+v", code, message, err)
		} else if util.IsExist(code, errs.WarnCodeListByRPC) {
			logz.Warnf(ctx, "rpc error code: %d，error message: %s，error: %+v", code, message, err)
		} else {
			logc.Errorf(ctx, "rpc error code: %d，error message: %s，error: %+v", code, message, err)
		}
	} else {
		// gateway业务错误
		code, message = errs.GetErrCodeAndMessage(err)
		_, file, line, _ := runtime.Caller(1)
		// 打印不同错误级别的日志
		if util.IsExist(code, errs.InfoCodeListByGateway) {
			logz.Infof(ctx, "报错地址：%s:%d，错误链路：%s，错误码：%v，错误原因：%s", file, line, err.Error(), code, message)
		} else if util.IsExist(code, errs.WarnCodeListByGateway) {
			logz.Warnf(ctx, "报错地址：%s:%d，错误链路：%s，错误码：%v，错误原因：%s", file, line, err.Error(), code, message)
		} else {
			logc.Errorf(ctx, "报错地址：%s:%d，错误链路：%s，错误码：%v，错误原因：%s", file, line, err.Error(), code, message)
		}
	}

	return code, message
}
