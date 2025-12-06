package course

import (
	"muse-admin/pkg/errs"
	"muse-admin/pkg/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"muse-admin/internal/logic/kingclub/course"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
)

func GetCourseChapterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetCourseChapterReq
		if err := httpx.Parse(r, &req); err != nil {
			response.ErrorCtx(r.Context(), w, errs.WithCode(err, errs.ErrCodeParamsAbnormal), nil)
			return
		}

		l := course.NewGetCourseChapter(r.Context(), svcCtx)
		resp, err := l.GetCourseChapter(&req)
		if err != nil {
			response.ErrorCtx(r.Context(), w, err, resp)
		} else {
			response.SuccessCtx(r.Context(), w, resp)
		}
	}
}
