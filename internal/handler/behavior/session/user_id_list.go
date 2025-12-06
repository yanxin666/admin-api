package session

import (
	"muse-admin/pkg/response"
	"net/http"

	"muse-admin/internal/logic/behavior/session"
	"muse-admin/internal/svc"
)

func UserIdListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := session.NewUserIdList(r.Context(), svcCtx)
		resp, err := l.UserIdList()
		if err != nil {
			response.ErrorCtx(r.Context(), w, err, resp)
		} else {
			response.SuccessCtx(r.Context(), w, resp)
		}
	}
}
