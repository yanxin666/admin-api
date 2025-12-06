package user

import (
	"muse-admin/internal/logic/user"
	"muse-admin/internal/svc"
	"muse-admin/pkg/response"
	"net/http"
)

func GetUserProfileInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewGetUserProfileInfoLogic(r.Context(), svcCtx)
		resp, err := l.GetUserProfileInfo()
		if err != nil {
			response.ErrorCtx(r.Context(), w, err, nil)
			return
		}

		response.SuccessCtx(r.Context(), w, resp)
	}
}
