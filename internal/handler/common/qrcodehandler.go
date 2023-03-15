package common

import (
	"chatgpt-tools/common/errorx"
	"chatgpt-tools/internal/types"
	"net/http"

	"chatgpt-tools/internal/logic/common"
	"chatgpt-tools/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func QrcodeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QrCodeRequest
		if err := httpx.Parse(r, &req); err != nil {
			errorx.Error(w, err.Error())
			return
		}

		l := common.NewQrcodeLogic(r.Context(), svcCtx)
		resp, err := l.Qrcode(req)
		if err != nil {
			errorx.Error(w, err.Error())
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
