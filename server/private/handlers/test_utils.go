package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
)

func makeRequest(
	method, url string, handle http.HandlerFunc, body io.Reader,
) *http.Response {

	req := httptest.NewRequest(method, url, body)
	if body != nil && method == http.MethodPost {
		ct := "application/x-www-form-urlencoded"
		req.Header.Set("Content-Type", ct)
	}

	w := httptest.NewRecorder()
	handle(w, req)

	return w.Result()
}
