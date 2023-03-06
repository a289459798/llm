package middleware

import (
	"chatgpt-tools/common/errorx"
	"chatgpt-tools/internal/config"
	"chatgpt-tools/model"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

type AuthAndUseMiddleware struct {
	C  config.Config
	DB *gorm.DB
}

func NewAuthAndUseMiddleware(c config.Config, db *gorm.DB) *AuthAndUseMiddleware {
	return &AuthAndUseMiddleware{
		C:  c,
		DB: db,
	}
}

func (m *AuthAndUseMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		authorization := r.Header.Get("authorization")
		uid := r.Context().Value("uid")
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

			str := fmt.Sprintf("%s-%s-%s", uid, timestamp, authorization)
			md5sum := md5.Sum([]byte(str))
			newSign := hex.EncodeToString(md5sum[:16])

			if newSign != sign {
				errorx.BadRequest(w, "参数错误")
				return
			}
		}

		uid2, _ := uid.(json.Number).Int64()
		amount := model.NewAccount(m.DB).GetAccount(uint32(uid2), time.Now())
		if amount.ChatAmount <= amount.ChatUse {
			errorx.Error(w, "次数已用完")
			return
		}

		next(w, r)
	}
}
