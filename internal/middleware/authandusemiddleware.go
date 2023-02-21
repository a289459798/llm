package middleware

import (
	"chatgpt-tools/common/errorx"
	"chatgpt-tools/internal/config"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type AuthAndUseMiddleware struct {
	C config.Config
}

func NewAuthAndUseMiddleware(c config.Config) *AuthAndUseMiddleware {
	return &AuthAndUseMiddleware{
		C: c,
	}
}

func (m *AuthAndUseMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if m.C.Mode == "pro" {
			timestamp := r.Header.Get("timestamp")
			sign := r.Header.Get("sign")
			if timestamp == "" {
				errorx.BadRequest(w, "参数错误")
				return
			}
			if sign == "" {
				errorx.BadRequest(w, "参数错误")
				return
			}

			t, _ := strconv.Atoi(timestamp)
			t1 := int64(t / 1000)

			if time.Now().Unix()-t1 > 60 {
				errorx.BadRequest(w, "参数错误")
				return
			}

			authorization := r.Header.Get("authorization")
			uid := r.Context().Value("uid")

			str := fmt.Sprintf("%s-%s-%s", uid, timestamp, authorization)
			md5sum := md5.Sum([]byte(str))
			newSign := hex.EncodeToString(md5sum[:16])

			if newSign != sign {
				errorx.BadRequest(w, "参数错误")
				return
			}
		}

		next(w, r)
	}
}
