package brain

import (
	"chatgpt-tools/common/errorx"
	"chatgpt-tools/internal/types"
	"net/http"

	"chatgpt-tools/internal/logic/chat/brain"
	"chatgpt-tools/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ChatHistoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChatHistoryRequest
		if err := httpx.Parse(r, &req); err != nil {
			errorx.Error(w, err.Error())
			return
		}

		l := brain.NewChatHistoryLogic(r.Context(), svcCtx)
		resp, err := l.ChatHistory(req)
		if err != nil {
			errorx.Error(w, err.Error())
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
