package benefit

import (
	"muse-admin/pkg/response"
	"net/http"

	"muse-admin/internal/logic/benefit"
	"muse-admin/internal/svc"
)

func GroupListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := benefit.NewGroupList(r.Context(), svcCtx)
		resp, err := l.GroupList()
		if err != nil {
			response.ErrorCtx(r.Context(), w, err, resp)
		} else {
			response.SuccessCtx(r.Context(), w, resp)
		}
	}
}
