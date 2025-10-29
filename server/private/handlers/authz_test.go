package handlers

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthzHandlerShouldRenderFormsOnGet(t *testing.T) {
	t.Parallel()

	type testCase struct {
		url      string
		router   http.HandlerFunc
		contains string
	}

	testCases := []testCase{{
		url:      "/signup",
		router:   authzRouter.GetSignup,
		contains: "Signup",
	}, {
		url:      "/login",
		router:   authzRouter.GetLogin,
		contains: "Login",
	}}
	for _, tc := range testCases {
		res := makeRequest(http.MethodGet, tc.url, tc.router)
		assert.Equal(t, "text/html", res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusOK, res.StatusCode)

		data, err := io.ReadAll(res.Body)
		body := string(data)
		t.Cleanup(func() { res.Body.Close() })
		assert.NoError(t, err)
		assert.Contains(t, body, tc.contains)
		assert.Contains(t, body, "</form>")
	}
}
