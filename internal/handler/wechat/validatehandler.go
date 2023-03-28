package wechat

import (
	"chatgpt-tools/internal/types"
	"net/http"

	"chatgpt-tools/internal/logic/wechat"
	"chatgpt-tools/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ValidateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WechatValidateRequest
		l := wechat.NewValidateLogic(r.Context(), svcCtx)
		resp, err := l.Validate(req, w, r)
		if err != nil {
			httpx.Error(w, err)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(resp))
		}
	}
}
