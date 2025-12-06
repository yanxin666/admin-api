package dept

import (
	"muse-admin/internal/logic/sys/dept"
	"muse-admin/internal/svc"
	"muse-admin/pkg/response"
	"net/http"
)

func GetSysDeptListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := dept.NewGetSysDeptListLogic(r.Context(), svcCtx)
		resp, err := l.GetSysDeptList()
		if err != nil {
			response.ErrorCtx(r.Context(), w, err, nil)
			return
		}

		response.SuccessCtx(r.Context(), w, resp)
	}
}
