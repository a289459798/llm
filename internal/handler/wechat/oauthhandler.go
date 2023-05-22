package wechat

import (
	"net/http"

	"chatgpt-tools/internal/logic/wechat"
	"chatgpt-tools/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func OauthHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := wechat.NewOauthLogic(r.Context(), svcCtx)
		resp, err := l.Oauth(r)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
