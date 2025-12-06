package role

import (
	"muse-admin/internal/logic/sys/role"
	"muse-admin/internal/svc"
	"muse-admin/pkg/response"
	"net/http"
)

func GetSysRoleListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := role.NewGetSysRoleListLogic(r.Context(), svcCtx)
		resp, err := l.GetSysRoleList()
		if err != nil {
			response.ErrorCtx(r.Context(), w, err, nil)
			return
		}

		response.SuccessCtx(r.Context(), w, resp)
	}
}
