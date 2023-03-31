package middleware

import "net/http"

type CronMiddleMiddleware struct {
}

func NewCronMiddleMiddleware() *CronMiddleMiddleware {
	return &CronMiddleMiddleware{}
}

func (m *CronMiddleMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r)
	}
}
