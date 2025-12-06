package tools

import (
	"context"
	"encoding/json"
	"muse-admin/internal/define"
	"muse-admin/internal/types"
)

// GetUserIdByCtx 上下文中获取用户ID
func GetUserIdByCtx(ctx context.Context) int64 {
	if jsonUid, ok := ctx.Value(define.SysJwtUserId).(json.Number); ok {
		if int64Uid, err := jsonUid.Int64(); err == nil {
			return int64Uid
		}
	}
	return 0
}

// GetUserInfoByCtx 上下文中获取用户信息
func GetUserInfoByCtx(ctx context.Context) types.UserInfoCtx {
	userCtx := ctx.Value(define.SysUserInfoCtx)
	switch expr := userCtx.(type) {
	case types.UserInfoCtx:
		return expr
	default:
		return types.UserInfoCtx{}
	}
}

// GetRequestParamsByCtx 根据上下文获取接口请求参数
func GetRequestParamsByCtx(ctx context.Context) map[string]any {
	params := ctx.Value(define.CtxRequestParams)
	switch expr := params.(type) {
	case map[string]any:
		return expr
	default:
		return nil
	}
}

// GetRequestLogsByCtx 根据上下文获取请求日志中的数据
func GetRequestLogsByCtx(ctx context.Context) types.MiddlewareApiRequestLog {
	params := ctx.Value(define.CtxRequestLogs)
	switch expr := params.(type) {
	case types.MiddlewareApiRequestLog:
		return expr
	default:
		return types.MiddlewareApiRequestLog{}
	}
}
