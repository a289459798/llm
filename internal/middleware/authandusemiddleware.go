package middleware

import (
	"net/http"
)

type AuthAndUseMiddleware struct {
}

func NewAuthAndUseMiddleware() *AuthAndUseMiddleware {
	return &AuthAndUseMiddleware{}
}

func (m *AuthAndUseMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO generate middleware implement function, delete after code implementation

		// Passthrough to next handler if need
		next(w, r)
	}
}
