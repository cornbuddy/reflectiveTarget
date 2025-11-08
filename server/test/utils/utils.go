package utils

import (
	"io"
	"net/http"
	"net/http/httptest"
)

func MakeRequest(
	contentType, method string,
	handle http.HandlerFunc, body io.Reader,
) *http.Response {

	req := httptest.NewRequest(method, "/", body)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	w := httptest.NewRecorder()
	handle(w, req)

	return w.Result()
}
