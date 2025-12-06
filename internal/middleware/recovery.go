package middleware

import (
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/rest/httpx"
	"muse-admin/internal/config"
	"net/http"
	"runtime/debug"
)

type RecoveryMiddleware struct {
	c config.Config
}

func NewRecoveryMiddleware(c config.Config) *RecoveryMiddleware {
	return &RecoveryMiddleware{
		c: c,
	}
}

func (m *RecoveryMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// var cls *tencent.Cls
		defer func() {
			// defer func() {
			// 	// 腾讯云日志上传
			// 	if cls != nil {
			// 		defer cls.Flush(2000)
			// 	}
			// }()
			if err := recover(); err != nil {
				logc.Errorf(r.Context(), "server is panic, stark：%s, error：%s", string(debug.Stack()), err)
				httpx.ErrorCtx(r.Context(), w, errors.New("程序异常"))
			}
		}()

		// // CLS日志服务
		// cls = auto_log.SetTencentLoggerProducer(config.AutoCls{
		// 	SecretId:  m.c.TencentCloud.SecretId,
		// 	SecretKey: m.c.TencentCloud.SecretKey,
		// 	TopicID:   m.c.Cls.TopicID,
		// 	Endpoint:  m.c.Cls.Endpoint,
		// 	Mode:      m.c.Mode,
		// }, r.Context())
		next(w, r)
	}
}
