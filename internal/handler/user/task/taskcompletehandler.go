package task

import (
	"chatgpt-tools/common/errorx"
	"net/http"

	"chatgpt-tools/internal/logic/user/task"
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func TaskCompleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TaskRequest
		if err := httpx.Parse(r, &req); err != nil {
			errorx.Error(w, err.Error())
			return
		}

		l := task.NewTaskCompleteLogic(r.Context(), svcCtx)
		resp, err := l.TaskComplete(&req, r)
		if err != nil {
			errorx.Error(w, err.Error())
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
