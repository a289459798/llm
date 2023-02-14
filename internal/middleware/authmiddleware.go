package middleware

import (
	"chatgpt-tools/common/errorx"
	"chatgpt-tools/internal/config"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
)

type AuthMiddleware struct {
	C config.Config
}

func NewAuthMiddleware(c config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		C: c,
	}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("authorization")
		if authorization == "" {
			errorx.Unauthorized(w, "no login")
			return
		}
		token, err := jwt.Parse(authorization, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(m.C.JwtSecret), nil
		})
		if err != nil {
			errorx.Unauthorized(w, err.Error())
			return
		}
		jwtData := token.Claims.(jwt.MapClaims)
		ctx := context.WithValue(r.Context(), "uid", int(jwtData["uid"].(float64)))
		next(w, r.WithContext(ctx))
	}
}
