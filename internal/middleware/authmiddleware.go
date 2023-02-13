package middleware

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
)

type AuthMiddleware struct {
}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("authorization")
		if authorization == "" {
			//response.Unauthorized(w, "no login")
			return
		}
		token, err := jwt.Parse(authorization, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("secret"), nil
		})
		if err != nil {
			//response.Unauthorized(w, err.Error())
			return
		}

		jwtData := token.Claims.(jwt.MapClaims)

		//val := m.RedisConn.HGet(r.Context(), "PetDoctorCloudMarketSalesmanCode", strconv.Itoa(int(jwtData["uid"].(float64))))
		//val2, _ := strconv.ParseFloat(val.Val(), 64)
		//if val2 != jwtData["code"] {
		//	response.Unauthorized(w, "登录过期")
		//	return
		//}
		//
		//salesman, _ := m.PetSalesmanModel.FindOneById(r.Context(), uint64(jwtData["uid"].(float64)))
		//if salesman == nil {
		//	response.Unauthorized(w, "用户不存在")
		//	return
		//}
		c := context.WithValue(r.Context(), "uid", int(jwtData["uid"].(float64)))
		ctx := context.WithValue(c, "user", "")
		next(w, r.WithContext(ctx))
	}
}
