package brain

import (
	"chatgpt-tools/common/errorx"
	"net/http"

	"chatgpt-tools/internal/logic/chat/brain"
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ChatHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChatRequest
		if err := httpx.Parse(r, &req); err != nil {
			errorx.Error(w, err.Error())
			return
		}

		l := brain.NewChatLogic(r.Context(), svcCtx)
		_, err := l.Chat(&req, w, r)
		if err != nil {
			errorx.Error(w, err.Error())
		}
	}
}
