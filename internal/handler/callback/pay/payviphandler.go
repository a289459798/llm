package pay

import (
	"net/http"

	"chatgpt-tools/internal/logic/callback/pay"
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func PayVipHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PayRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := pay.NewPayVipLogic(r.Context(), svcCtx)
		resp, err := l.PayVip(&req, r)
		if err != nil {
			httpx.Error(w, err)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(resp.Data))
		}
	}
}
