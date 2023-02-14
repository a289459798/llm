package divination

import (
	"chatgpt-tools/common/errorx"
	"net/http"

	"chatgpt-tools/internal/logic/divination"
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GsqimingHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GSQiMingRequest
		if err := httpx.Parse(r, &req); err != nil {
			errorx.Error(w, err.Error())
			return
		}

		l := divination.NewGsqimingLogic(r.Context(), svcCtx)
		resp, err := l.Gsqiming(&req, w)
		if err != nil {
			errorx.Error(w, err.Error())
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
