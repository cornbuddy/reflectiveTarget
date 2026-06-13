package utils

import (
	"net/http"
)

// writes response body according to http status
func HttpError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
