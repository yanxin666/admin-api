package course

import (
	"muse-admin/pkg/response"
	"net/http"

	"muse-admin/internal/logic/kingclub/course"
	"muse-admin/internal/svc"
)

func GetCourseCommonDataHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := course.NewGetCourseCommonData(r.Context(), svcCtx)
		resp, err := l.GetCourseCommonData()
		if err != nil {
			response.ErrorCtx(r.Context(), w, err, resp)
		} else {
			response.SuccessCtx(r.Context(), w, resp)
		}
	}
}
