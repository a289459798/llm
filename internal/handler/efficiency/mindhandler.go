package efficiency

import (
	"chatgpt-tools/common/errorx"
	"net/http"

	"chatgpt-tools/internal/logic/efficiency"
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func MindHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MindRequest
		if err := httpx.Parse(r, &req); err != nil {
			errorx.Error(w, err.Error())
			return
		}

		l := efficiency.NewMindLogic(r.Context(), svcCtx)
		resp, err := l.Mind(&req)
		if err != nil {
			errorx.Error(w, err.Error())
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
