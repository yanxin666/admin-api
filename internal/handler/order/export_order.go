package order

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"muse-admin/internal/logic/order"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
)

func ExportOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetOrderListReq
		if err := httpx.Parse(r, &req); err != nil {
			// todo: replace it with the return function you need to call
			httpx.ErrorCtx(r.Context(), w, err)
			// response.ErrorCtx(r.Context(), w, errs.WithCode(err, errs.ErrCodeParamsAbnormal), nil)
			return
		}

		l := order.NewExportOrder(r.Context(), svcCtx)
		resp, err := l.ExportOrder(&req)
		if err != nil {
			// todo: replace it with the return function you need to call
			httpx.ErrorCtx(r.Context(), w, err)
			// response.ErrorCtx(r.Context(), w, err, resp)
		} else {
			// todo: replace it with the return function you need to call
			httpx.OkJsonCtx(r.Context(), w, resp)
			// response.SuccessCtx(r.Context(), w, resp)
		}
	}
}
