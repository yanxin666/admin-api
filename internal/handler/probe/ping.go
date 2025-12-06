package probe

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"muse-admin/internal/logic/probe"
	"muse-admin/internal/svc"
)

// PingHandler ping服务探活
func PingHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := probe.NewPing(r.Context(), svcCtx)
		err := l.Ping()

		// _, err = svcCtx.MQProducer.SendSync(r.Context(), mqdef.TopicTest, "this is my wang")

		// 不要修改这里的返回方法为response
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, "success")
		}
	}
}
