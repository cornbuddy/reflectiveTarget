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
) (*http.Response, string, error) {

	req := httptest.NewRequest(method, url, body)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	return makeRequest(req, handle)
}

func MakeRequest(
	contentType, method, url string,
	handle http.HandlerFunc, body io.Reader,
) (*http.Response, string, error) {

	req := httptest.NewRequest(method, url, body)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	return makeRequest(req, handle)
}

func makeRequest(
	req *http.Request, handle http.HandlerFunc,
) (*http.Response, string, error) {

	w := httptest.NewRecorder()
	handle(w, req)

	r := w.Result()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, "", err
	}

	defer r.Body.Close()

	return r, string(body), nil
}
