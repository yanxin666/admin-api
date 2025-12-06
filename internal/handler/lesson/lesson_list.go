package lesson

import (
	"muse-admin/pkg/errs"
	"muse-admin/pkg/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"muse-admin/internal/logic/lesson"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
)

func LessonListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LessonListReq
		if err := httpx.Parse(r, &req); err != nil {
			response.ErrorCtx(r.Context(), w, errs.WithCode(err, errs.ErrCodeParamsAbnormal), nil)
			return
		}

		l := lesson.NewLessonList(r.Context(), svcCtx)
		resp, err := l.LessonList(&req)
		if err != nil {
			response.ErrorCtx(r.Context(), w, err, resp)
		} else {
			response.SuccessCtx(r.Context(), w, resp)
		}
	}
}
