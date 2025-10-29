package handlers

import (
	"net/http"
	"net/http/httptest"
)

func makeRequest(method, url string, handle http.HandlerFunc) *http.Response {
	req := httptest.NewRequest(method, url, nil)
	w := httptest.NewRecorder()
	handle(w, req)

	return w.Result()
}
