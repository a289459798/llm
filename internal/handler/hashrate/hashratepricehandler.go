package hashrate

import (
	"chatgpt-tools/common/errorx"
	"net/http"

	"chatgpt-tools/internal/logic/hashrate"
	"chatgpt-tools/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func HashRatePriceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := hashrate.NewHashRatePriceLogic(r.Context(), svcCtx)
		resp, err := l.HashRatePrice()
		if err != nil {
			errorx.Error(w, err.Error())
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
