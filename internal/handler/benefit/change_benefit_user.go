package benefit

import (
	"net/http"

	"muse-admin/internal/logic/benefit"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"
	"muse-admin/pkg/response"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ChangeBenefitUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChangeBenefitUserReq
		if err := httpx.Parse(r, &req); err != nil {
			response.ErrorCtx(r.Context(), w, errs.WithCode(err, errs.ErrCodeParamsAbnormal), nil)
			return
		}

		l := benefit.NewChangeBenefitUser(r.Context(), svcCtx)
		resp, err := l.ChangeBenefitUser(&req)
		if err != nil {
			response.ErrorCtx(r.Context(), w, err, resp)
		} else {
			response.SuccessCtx(r.Context(), w, resp)
		}
	}
}
