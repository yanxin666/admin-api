package menu

import (
	"muse-admin/internal/logic/sys/menu"
	"muse-admin/internal/svc"
	"muse-admin/pkg/response"
	"net/http"
)

func GetSysPermMenuListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := menu.NewGetSysPermMenuListLogic(r.Context(), svcCtx)
		resp, err := l.GetSysPermMenuList()
		if err != nil {
			response.ErrorCtx(r.Context(), w, err, nil)
			return
		}

		response.SuccessCtx(r.Context(), w, resp)
	}
}
