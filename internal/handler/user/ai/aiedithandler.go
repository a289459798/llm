package ai

import (
	"chatgpt-tools/common/errorx"
	"net/http"

	"chatgpt-tools/internal/logic/user/ai"
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AiEditHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AIEditRequest
		if err := httpx.Parse(r, &req); err != nil {
			errorx.Error(w, err.Error())
			return
		}

		l := ai.NewAiEditLogic(r.Context(), svcCtx)
		uploadedFiles := r.MultipartForm.File
		resp, err := l.AiEdit(&req, uploadedFiles)
		if err != nil {
			errorx.Error(w, err.Error())
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
