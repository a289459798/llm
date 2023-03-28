package wechat

import (
	"chatgpt-tools/internal/types"
	"net/http"

	"chatgpt-tools/internal/logic/wechat"
	"chatgpt-tools/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func EventHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WechatValidateRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		l := wechat.NewEventLogic(r.Context(), svcCtx)
		resp, err := l.Event(req, r, w)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
