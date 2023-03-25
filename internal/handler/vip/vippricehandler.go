package vip

import (
	"net/http"

	"chatgpt-tools/internal/logic/vip"
	"chatgpt-tools/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func VipPriceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := vip.NewVipPriceLogic(r.Context(), svcCtx)
		resp, err := l.VipPrice()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
