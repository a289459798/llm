package user

import (
	"chatgpt-tools/common/errorx"
	"chatgpt-tools/internal/types"
	"net/http"

	"chatgpt-tools/internal/logic/user"
	"chatgpt-tools/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func LoginWXQrcodeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WxqrcodeRequest
		if err := httpx.Parse(r, &req); err != nil {
			errorx.Error(w, err.Error())
			return
		}
		l := user.NewLoginWXQrcodeLogic(r.Context(), svcCtx)
		resp, err := l.LoginWXQrcode(req, r)
		if err != nil {
			errorx.Error(w, err.Error())
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
