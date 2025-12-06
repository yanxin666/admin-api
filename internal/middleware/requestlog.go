package middleware

import (
	"bytes"
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/third/tencent"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/jsonx"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"io"
	"muse-admin/internal/config"
	"muse-admin/internal/define"
	"muse-admin/internal/types"
	"net/http"
)

// RequestLogMiddleware 请求参数中间件
type RequestLogMiddleware struct {
	c config.Config
}

func NewRequestLogMiddleware(c config.Config) *RequestLogMiddleware {
	return &RequestLogMiddleware{c: c}
}

func (m *RequestLogMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取请求参数
		requestData := make(map[string]any, 30)
		if r.Method == http.MethodGet {
			query := r.URL.Query()
			for k := range query {
				requestData[k] = query.Get(k)
			}
		} else if r.Method == http.MethodPost {
			if r.Header.Get("Content-Type") == "application/json" {
				// 读取原始请求体
				bodyBytes, err := io.ReadAll(r.Body)

				// 关闭原始请求体
				defer func() {
					err = r.Body.Close()
					if err != nil {
						logc.Errorf(r.Context(), "关闭原始请求体失败，原因：%s", err)
					}
				}()

				// 读取请求体内容
				var params any
				reader := io.LimitReader(bytes.NewReader(bodyBytes), 8<<20)
				err = jsonx.UnmarshalFromReader(reader, &params)
				if err == nil {
					for k, v := range params.(map[string]any) {
						requestData[k] = v
					}
				}

				// 重新设置请求体为已经读取的内容
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			} else {
				requestData, _ = httpx.GetFormValues(r)
			}
		}

		// 设置日志自定义内容
		logFields := []logx.LogField{
			logx.Field(define.LogFieldsType.Path, r.URL.Path),
		}
		ctx := logx.ContextWithFields(r.Context(), logFields...)
		r = r.WithContext(ctx)

		// 记录请求日志
		requestParam := types.MiddlewareApiRequestLog{
			Host:     r.Host,
			ClientIP: r.RemoteAddr,
			Schema:   r.Proto,
			Header:   r.Header,
			URL:      r.URL.String(),
			Path:     r.URL.Path,
			Method:   r.Method,
			Params:   requestData,
			UserId:   GetUserId(ctx),
		}
		requestParamJson, _ := json.Marshal(&requestParam)
		logc.Info(r.Context(), "[ApiRequestIn] ", string(requestParamJson))

		// 设置请求参数到上下文
		ctx = context.WithValue(r.Context(), define.CtxRequestParams, requestData)
		r = r.WithContext(ctx)

		// 设置请求参数到上下文
		ctx = context.WithValue(ctx, define.CtxRequestLogs, requestParam)
		r = r.WithContext(ctx)

		next(w, r)
	}
}

// tencentCls 腾讯云日志服务
func tencentCls(c config.Config, ctx context.Context) *tencent.Cls {
	var cls *tencent.Cls
	if c.Cls.TopicID != "" && c.Cls.Endpoint != "" && c.Mode != "dev" {
		// 腾讯云日志
		cls = tencent.NewCls(c.TencentCloud.SecretId, c.TencentCloud.SecretKey, c.Cls.TopicID, c.Cls.Endpoint, 20)
		err := cls.GetProducerClient(context.Background())
		if err != nil {
			logc.Errorf(ctx, "初始化腾讯云日志服务异常：%s", err)
		} else {
			logx.SetWriter(logx.NewWriter(cls))
		}
	}

	return cls
}
