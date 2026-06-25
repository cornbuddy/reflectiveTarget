package utils

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
)

func MakeRequestWithCookies(
	contentType, method, url string,
	handle http.HandlerFunc, body io.Reader,
	cookies ...*http.Cookie,
) (*http.Response, string, error) {
	//nolint:noctx
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
	//nolint:noctx
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
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, "", err
	}

	if err := r.Body.Close(); err != nil {
		return nil, "", err
	}

	r.Body = io.NopCloser(bytes.NewBuffer(raw))

	return r, string(raw), nil
}
