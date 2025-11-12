package utils

import (
	"io"
	"net/http"
	"net/http/httptest"
)

func MakeRequestWithCookies(
	contentType, method, url string,
	handle http.HandlerFunc, body io.Reader,
	cookies ...*http.Cookie,
) *http.Response {

	req := httptest.NewRequest(method, url, body)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	return MakeRequest(contentType, method, url, handle, body)
}

func MakeRequest(
	contentType, method, url string,
	handle http.HandlerFunc, body io.Reader,
) *http.Response {

	req := httptest.NewRequest(method, url, body)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	return makeRequest(req, handle)
}

func makeRequest(req *http.Request, handle http.HandlerFunc) *http.Response {
	w := httptest.NewRecorder()
	handle(w, req)

	return w.Result()
}
