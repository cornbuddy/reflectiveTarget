package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeRequest(
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

func assertSessionCookieIsSet(t *testing.T, resp *http.Response) {
	cookies := resp.Cookies()
	require.NotEmpty(t, cookies)

	var sessionCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == SessionCookieName {
			sessionCookie = cookie
			break
		}
	}

	require.NotNil(t, sessionCookie)

	// substracting 1 second because time.Now add ms to the time, while
	// cookie doesn't count that strictly
	month := time.Now().AddDate(0, 1, 0).Add(-1 * time.Second).UTC()
	assert.True(t, sessionCookie.Expires.After(month))
	assert.NoError(t, uuid.Validate(sessionCookie.Value))
}
