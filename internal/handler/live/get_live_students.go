package live

import (
	"muse-admin/pkg/errs"
	"net/http"

	"muse-admin/internal/logic/live"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"muse-admin/pkg/response"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetLiveStudentsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetLiveStudentsReq
		if err := httpx.Parse(r, &req); err != nil {
			response.ErrorCtx(r.Context(), w, errs.WithCode(err, errs.ErrCodeParamsAbnormal), nil)
			return
		}

		l := live.NewGetLiveStudents(r.Context(), svcCtx)
		resp, err := l.GetLiveStudents(&req)
		if err != nil {
			response.ErrorCtx(r.Context(), w, err, resp)
		} else {
			response.SuccessCtx(r.Context(), w, resp)
		}
	}
}
