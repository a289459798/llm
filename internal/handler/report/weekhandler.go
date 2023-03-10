package report

import (
	"chatgpt-tools/common/errorx"
	"chatgpt-tools/internal/logic/report"
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func WeekHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ReportRequest
		if err := httpx.Parse(r, &req); err != nil {
			errorx.Error(w, err.Error())
			return
		}

		l := report.NewWeekLogic(r.Context(), svcCtx)
		_, err := l.Week(&req, w)
		if err != nil {
			errorx.Error(w, err.Error())
		}
	}
}
