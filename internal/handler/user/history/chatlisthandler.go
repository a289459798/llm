package history

import (
	"chatgpt-tools/common/errorx"
	"chatgpt-tools/internal/types"
	"net/http"

	"chatgpt-tools/internal/logic/user/history"
	"chatgpt-tools/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ChatListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PageRequest
		if err := httpx.Parse(r, &req); err != nil {
			errorx.Error(w, err.Error())
			return
		}

		l := history.NewChatListLogic(r.Context(), svcCtx)
		resp, err := l.ChatList(req)
		if err != nil {
			errorx.Error(w, err.Error())
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
