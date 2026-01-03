package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

func TestMakeHttpHandler(t *testing.T) {
	t.Parallel()

	type testCase struct {
		url        string
		method     string
		statusCode int
	}

	testCases := []testCase{{
		url:        "/",
		method:     http.MethodGet,
		statusCode: http.StatusOK,
	}, {
		url:        "/kek",
		method:     http.MethodGet,
		statusCode: http.StatusNotFound,
	}, {
		url:        "/login",
		method:     http.MethodGet,
		statusCode: http.StatusOK,
	}, {
		url:        "/login",
		method:     http.MethodPost,
		statusCode: http.StatusUnauthorized,
	}, {
		url:        "/signup",
		method:     http.MethodGet,
		statusCode: http.StatusOK,
	}, {
		url:        "/signup",
		method:     http.MethodPost,
		statusCode: http.StatusBadRequest,
	}, {
		url:        "/targets",
		method:     http.MethodGet,
		statusCode: http.StatusForbidden,
	}, {
		url:        "/targets",
		method:     http.MethodPut,
		statusCode: http.StatusForbidden,
	}, {
		url:        "/targets",
		method:     http.MethodPost,
		statusCode: http.StatusForbidden,
	}, {
		url:        "/api/health",
		method:     http.MethodGet,
		statusCode: http.StatusOK,
	}, {
		url:        "/api/target/1/shots",
		method:     http.MethodGet,
		statusCode: http.StatusNotFound,
	}, {
		url:        "/api/target/1/shots",
		method:     http.MethodPost,
		statusCode: http.StatusBadRequest,
	}}

	for _, tc := range testCases {
		resp := utils.MakeRequest("", tc.method, tc.url, router, nil)
		assertSessionCookieIsSet(t, resp)

		msg := fmt.Sprintf("%s %s", tc.method, tc.url)
		assert.Equal(t, tc.statusCode, resp.StatusCode, msg)
	}
}
