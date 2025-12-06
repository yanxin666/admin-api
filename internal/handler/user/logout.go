package user

import (
	"muse-admin/internal/logic/user"
	"muse-admin/internal/svc"
	"muse-admin/pkg/response"
	"net/http"
)

func LogoutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewLogoutLogic(r.Context(), svcCtx)
		err := l.Logout()
		if err != nil {
			response.ErrorCtx(r.Context(), w, err, nil)
			return
		}

		response.SuccessCtx(r.Context(), w, nil)
	}
}
