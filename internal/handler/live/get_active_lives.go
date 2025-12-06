package live

import (
	"net/http"

	"muse-admin/internal/logic/live"
	"muse-admin/internal/svc"
	"muse-admin/pkg/response"
)

func GetActiveLivesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := live.NewGetActiveLives(r.Context(), svcCtx)
		resp, err := l.GetActiveLives()
		if err != nil {
			response.ErrorCtx(r.Context(), w, err, resp)
		} else {
			response.SuccessCtx(r.Context(), w, resp)
		}
	}
}
