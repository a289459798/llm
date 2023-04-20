package pay

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"chatgpt-tools/internal/logic/callback/pay"
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func PayHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body []byte
		if r.Body != nil {
			body, _ = ioutil.ReadAll(r.Body)
			r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		}

		var req types.PayRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		l := pay.NewPayLogic(r.Context(), svcCtx)
		resp, err := l.Pay(&req, r)
		if err != nil {
			httpx.Error(w, err)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(resp.Data))
			return
		}
	}
}
