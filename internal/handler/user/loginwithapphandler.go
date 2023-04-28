package user

import (
	"chatgpt-tools/common/errorx"
	"net/http"

	"chatgpt-tools/internal/logic/user"
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func LoginWithAppHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginAppRequest
		if err := httpx.Parse(r, &req); err != nil {
			errorx.Error(w, err.Error())
			return
		}

		l := user.NewLoginWithAppLogic(r.Context(), svcCtx)
		resp, err := l.LoginWithApp(&req)
		if err != nil {
			errorx.Error(w, err.Error())
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
