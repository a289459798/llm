package image

import (
	"chatgpt-tools/common/errorx"
	"chatgpt-tools/internal/logic/image"
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func Pic2picHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Pic2picRequest
		if err := httpx.ParseForm(r, &req); err != nil {
			errorx.Error(w, err.Error())
			return
		}

		l := image.NewPic2picLogic(r.Context(), svcCtx)

		uploadedFiles := r.MultipartForm.File

		resp, err := l.Pic2pic(&req, uploadedFiles)
		if err != nil {
			errorx.Error(w, err.Error())
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
