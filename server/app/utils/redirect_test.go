package utils

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedirect(t *testing.T) {
	t.Parallel()

	type testCase struct {
		desc               string
		requestHeaders     http.Header
		url                string
		wantBody           string
		wantResponseHeader http.Header
		wantStatus         int
	}

	testCases := []testCase{{
		"should return http redirect on malformed htmx header",
		http.Header{htmxRequestKey: []string{"false"}},
		"/",
		"kek",
		http.Header{"Location": []string{"/"}},
		http.StatusSeeOther,
	}, {
		"should return htmx redirect on htmx header",
		http.Header{htmxRequestKey: []string{"true"}},
		"/",
		"kek",
		http.Header{"Hx-Redirect": []string{"/"}},
		http.StatusNoContent,
	}, {
		"should return http redirect on empty header",
		http.Header{},
		"/",
		"kek",
		http.Header{"Location": []string{"/"}},
		http.StatusSeeOther,
	}}

	for _, tc := range testCases {
		handle := func(w http.ResponseWriter, r *http.Request) {
			Redirect(w, r, tc.url, tc.wantBody)
		}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/kek", nil)
		r.Header = tc.requestHeaders
		handle(w, r)
		resp := w.Result()
		data, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		t.Cleanup(func() {
			require.NoError(t, resp.Body.Close())
		})

		assert.Equal(t, tc.wantBody, string(data), tc.desc)
		assert.Equal(t, tc.wantResponseHeader, resp.Header, tc.desc)
		assert.Equal(t, tc.wantStatus, resp.StatusCode, tc.desc)
	}
}
