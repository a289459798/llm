package code

import (
	"chatgpt-tools/common/errorx"
	"net/http"

	"chatgpt-tools/internal/logic/code"
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ExamHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ExamRequest
		if err := httpx.Parse(r, &req); err != nil {
			errorx.Error(w, err.Error())
			return
		}

		l := code.NewExamLogic(r.Context(), svcCtx)
		_, err := l.Exam(&req, w)
		if err != nil {
			errorx.Error(w, err.Error())
		}
	}
}
