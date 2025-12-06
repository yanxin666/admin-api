package dict

import (
	"muse-admin/internal/logic/config/dict"
	"muse-admin/internal/svc"
	"muse-admin/pkg/response"
	"net/http"
)

func GetConfigDictListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := dict.NewGetConfigDictListLogic(r.Context(), svcCtx)
		resp, err := l.GetConfigDictList()
		if err != nil {
			response.ErrorCtx(r.Context(), w, err, nil)
			return
		}

		response.SuccessCtx(r.Context(), w, resp)
	}
}
