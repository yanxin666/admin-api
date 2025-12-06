package user

import (
	"muse-admin/internal/logic/user"
	"muse-admin/internal/svc"
	"muse-admin/pkg/response"
	"net/http"
)

func GetUserInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewGetUserInfoLogic(r.Context(), svcCtx)
		resp, err := l.GetUserInfo()
		if err != nil {
			response.ErrorCtx(r.Context(), w, err, nil)
			return
		}

		response.SuccessCtx(r.Context(), w, resp)
	}
}
