package wechat

import (
	"net/http"

	"chatgpt-tools/internal/logic/wechat"
	"chatgpt-tools/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SetMenuHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := wechat.NewSetMenuLogic(r.Context(), svcCtx)
		resp, err := l.SetMenu()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
