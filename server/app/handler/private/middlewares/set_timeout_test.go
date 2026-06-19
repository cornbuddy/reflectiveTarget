package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSetTimeout(t *testing.T) {
	t.Parallel()

	const timeout = 10 * time.Millisecond

	type testCase struct {
		desc   string
		sleep  time.Duration
		status int
	}

	testCases := []testCase{{
		"should timeout if duration is too long",
		10 * timeout,
		http.StatusGatewayTimeout,
	}, {
		"should succeed if duration is met",
		timeout / 10,
		http.StatusOK,
	}}

	for _, tc := range testCases {
		stub := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(tc.sleep)
		})
		setTimeout := mw.SetTimeout(timeout)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		setTimeout(stub).ServeHTTP(w, r)

		assert.Equal(t, tc.status, w.Result().StatusCode, tc.desc)
	}
}
