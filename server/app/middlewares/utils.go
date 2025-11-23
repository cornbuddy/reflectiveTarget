package middlewares

import (
	"net/http"
)

func internalServerError(w http.ResponseWriter, msg string) {
	http.Error(w, msg, http.StatusInternalServerError)
}
