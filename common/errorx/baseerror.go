package errorx

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func Error(w http.ResponseWriter, err string) {
	switch err {
	case string(http.StatusNotFound):
		NoFound(w, "NotFound")
		break
	case string(http.StatusBadRequest):
		BadRequest(w, "BadRequest")
		break
	case string(http.StatusUnauthorized):
		Unauthorized(w, "Unauthorized")
		break
	default:
		httpx.WriteJson(w, http.StatusInternalServerError, map[string]string{"message": err})
		break
	}
}

func NoFound(w http.ResponseWriter, err string) {
	httpx.WriteJson(w, http.StatusNotFound, map[string]string{"message": err})
}

func Unauthorized(w http.ResponseWriter, err string) {
	httpx.WriteJson(w, http.StatusUnauthorized, map[string]string{"message": err})
}

func BadRequest(w http.ResponseWriter, err string) {
	httpx.WriteJson(w, http.StatusBadRequest, map[string]string{"message": err})
}

func Abort(w http.ResponseWriter, code int, err string) {
	httpx.WriteJson(w, code, map[string]string{"message": err})
}
