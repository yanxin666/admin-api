package response

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/trace"
	"github.com/zeromicro/go-zero/rest/httpx"
	"muse-admin/internal/tools"
	"muse-admin/pkg/errs"
	"net/http"
	"runtime"
)

type Body struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// EasyResponse 本框架初始化时的方法
func EasyResponse(w http.ResponseWriter, resp interface{}, err error) {
	var body Body
	if err != nil {
		body.Code = 0
		body.Msg = err.Error()
	} else {
		body.Code = 200
		body.Msg = "success"
		body.Data = resp
	}
	httpx.OkJson(w, body)
}

type Response struct {
	Code    errs.BasicCode `json:"code"`
	Message string         `json:"message"`
	Data    any            `json:"data"`
	TraceId string         `json:"traceId"`
}

type StreamEventResponse struct {
	Code    errs.BasicCode `json:"code"`
	Message string         `json:"message"`
	TraceId string         `json:"traceId"`
	Stop    bool           `json:"stop"`
	Error   bool           `json:"error"`
}

// SuccessCtx api接口json格式返回成功
func SuccessCtx(ctx context.Context, w http.ResponseWriter, v any) {
	response := Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    v,
		TraceId: trace.TraceIDFromContext(ctx),
	}
	httpx.OkJsonCtx(ctx, w, response)
}

// ErrorCtx api接口json格式返回失败
// err：考虑到结合框架中 error 处理函数按照上下文传递 errorCode 和 errorMessage ，并能自动填充到返回体的 code 和 message ，err 变量须由包 errs 中的方法生成
func ErrorCtx(ctx context.Context, w http.ResponseWriter, err error, v any) {
	code, message := tools.RpcErrChange(ctx, err)
	response := Response{
		Code:    code,
		Message: message,
		Data:    v,
		TraceId: trace.TraceIDFromContext(ctx),
	}
	httpx.WriteJsonCtx(ctx, w, http.StatusOK, response)
}

// SuccessSSECtx SSE返回数据
func SuccessSSECtx(ctx context.Context, w http.ResponseWriter, flusher http.Flusher, code errs.BasicCode, v any) {
	// SSE对象
	if flusher == nil {
		logc.Errorf(ctx, "SSE对象异常，error：%+v", flusher)
		return
	}

	// 设置头部信息
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	response := Response{
		Code:    code,
		Message: "success",
		Data:    v,
		TraceId: trace.TraceIDFromContext(ctx),
	}
	byt, err := json.Marshal(response)
	if err != nil {
		logc.Errorf(ctx, "SSE序列化异常，error：%s", err)
	}
	_, err = fmt.Fprintf(w, "data: %s\n\n", string(byt))
	if err != nil {
		logc.Errorf(ctx, "SSE success响应异常，error：%s", err)
	}
	flusher.Flush()
}

// EventSSECtx SSE返回事件
func EventSSECtx(ctx context.Context, w http.ResponseWriter, flusher http.Flusher, err error) {
	// SSE对象
	if flusher == nil {
		logc.Errorf(ctx, "SSE对象异常，error：%+v", flusher)
		return
	}

	// 设置头部信息
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// 获取错误码和错误信息
	code, message := errs.GetErrCodeAndMessage(err)
	response := StreamEventResponse{
		Code:    code,
		Message: message,
		TraceId: trace.TraceIDFromContext(ctx),
	}

	// 错误码分析
	_, file, line, _ := runtime.Caller(1)
	if code != errs.ErrCodeStreamStop && code != errs.ErrCodeFakeStreamStop {
		logc.Errorf(ctx, "接口报错：%s:%d，错误信息：%s", file, line, err.Error())
		response.Error = true
	} else {
		response.Stop = true
	}

	// 流式数据
	byt, err := json.Marshal(response)
	if err != nil {
		logc.Errorf(ctx, "SSE序列化异常，error：%s", err)
	}
	_, err = fmt.Fprintf(w, "data: %s\n\n", string(byt))
	if err != nil {
		logc.Errorf(ctx, "SSE event响应异常，error：%s， data：%v", err, response)
	}
	flusher.Flush()
}
