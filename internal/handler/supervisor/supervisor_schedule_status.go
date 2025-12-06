package supervisor

import (
	"muse-admin/pkg/errs"
	"muse-admin/pkg/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"muse-admin/internal/logic/supervisor"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
)

func SupervisorScheduleStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SupervisorScheduleStatusReq
		if err := httpx.Parse(r, &req); err != nil {
			response.ErrorCtx(r.Context(), w, errs.WithCode(err, errs.ErrCodeParamsAbnormal), nil)
			return
		}

		l := supervisor.NewSupervisorScheduleStatus(r.Context(), svcCtx)
		err := l.SupervisorScheduleStatus(&req)
		if err != nil {
			response.ErrorCtx(r.Context(), w, err, nil)
		} else {
			response.SuccessCtx(r.Context(), w, nil)
		}
	}
}
