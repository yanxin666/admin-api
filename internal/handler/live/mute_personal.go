package live

import (
	"muse-admin/pkg/errs"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"muse-admin/internal/logic/live"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"muse-admin/pkg/response"
)

func MutePersonalHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MutePersonalReq
		if err := httpx.Parse(r, &req); err != nil {
			response.ErrorCtx(r.Context(), w, errs.WithCode(err, errs.ErrCodeParamsAbnormal), nil)
			return
		}

		l := live.NewMutePersonal(r.Context(), svcCtx)
		err := l.MutePersonal(&req)
		if err != nil {
			response.ErrorCtx(r.Context(), w, err, nil)
		} else {
			response.SuccessCtx(r.Context(), w, nil)
		}
	}
}
