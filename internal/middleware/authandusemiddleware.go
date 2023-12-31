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
		appKey := r.Header.Get("App-Key")
		if appKey == "" {
			errorx.BadRequest(w, "参数错误 appKey")
			return
		}
		authorization := r.Header.Get("authorization")
		uid := r.Context().Value("uid")
		uid2, _ := uid.(json.Number).Int64()
		if m.C.Mode == "pro" {
			timestamp := r.Header.Get("timestamp")
			sign := r.Header.Get("sign")
			if timestamp == "" {
				errorx.BadRequest(w, "参数错误：timestamp")
				return
			}
			if sign == "" {
				errorx.BadRequest(w, "参数错误：sign")
				return
			}

			t, _ := strconv.Atoi(timestamp)
			t1 := int64(t / 1000)

			if time.Now().Unix()-t1 > 600 {
				errorx.BadRequest(w, "参数错误：时间过期")
				return
			}

			str := fmt.Sprintf("%s-%s-%s", uid, timestamp, authorization)
			md5sum := md5.Sum([]byte(str))
			newSign := hex.EncodeToString(md5sum[:16])

			if newSign != sign {
				errorData := model.Error{
					Uid:      uint32(uid2),
					Type:     "sign",
					Question: fmt.Sprintf("{\"timestamp\":\"%s\", \"authorization\":\"%s\",\"sign\":\"%s\",\"newSign\":\"%s\"}", timestamp, authorization, sign, newSign),
					Error:    "参数错误：验签失败",
				}
				errorData.Insert(m.DB)
				errorx.BadRequest(w, "参数错误：验签失败")
				return
			}

			requestLog := &model.RequestLog{}
			m.DB.Where("request_id = ?", newSign).First(&requestLog)
			if requestLog.ID > 0 {
				errorx.BadRequest(w, "重复请求")
				return
			}
			// 记录请求
			m.DB.Create(&model.RequestLog{RequestId: newSign})
		}

		amount := model.NewAccount(m.DB).GetAccount(uint32(uid2), time.Now())
		if amount.Amount <= 0 {
			errorx.Error(w, "次数已用完")
			return
		}

		next(w, r)

	}
}
