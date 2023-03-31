package crontab

import (
	"net/http"

	"chatgpt-tools/internal/logic/crontab"
	"chatgpt-tools/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func VipCheckHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := crontab.NewVipCheckLogic(r.Context(), svcCtx)
		resp, err := l.VipCheck()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
