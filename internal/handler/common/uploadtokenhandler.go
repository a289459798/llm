package common

import (
	"chatgpt-tools/common/errorx"
	"net/http"

	"chatgpt-tools/internal/logic/common"
	"chatgpt-tools/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UploadTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := common.NewUploadTokenLogic(r.Context(), svcCtx)
		resp, err := l.UploadToken()
		if err != nil {
			errorx.Error(w, err.Error())
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
