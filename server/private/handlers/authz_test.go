package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldRegisterNewUserWhenCredentialsAreValid(t *testing.T) {
	t.Parallel()

	type testCase struct {
		message    string
		statusCode int
		body       io.Reader
	}

	testCases := []testCase{{
		message:    "User successfully created",
		statusCode: http.StatusCreated,
		body: strings.NewReader(
			fmt.Sprintf("username=%s&password=%s", "kek", "kek"),
		),
	}, {
		message:    "User already exists",
		statusCode: http.StatusConflict,
		body: strings.NewReader(
			fmt.Sprintf("username=%s&password=%s", "kek", "kek"),
		),
	}}

	for _, tc := range testCases {
		route := authzRouter.PostSignup
		ct := "application/x-www-form-urlencoded"
		body := tc.body
		res := makeRequest(ct, http.MethodPost, "/signup", route, body)
		assert.Equal(t, tc.statusCode, res.StatusCode)

		data, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		t.Cleanup(func() { res.Body.Close() })

		gotBody := string(data)
		assert.Contains(t, gotBody, tc.message)
	}
}

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
		res := makeRequest("", http.MethodGet, tc.url, tc.router, nil)
		ct := "text/html; charset=utf-8"
		assert.Equal(t, ct, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusOK, res.StatusCode)

		data, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		t.Cleanup(func() { res.Body.Close() })

		body := string(data)
		assert.Contains(t, body, tc.contains)
		assert.Contains(t, body, "</form>")
	}
}
