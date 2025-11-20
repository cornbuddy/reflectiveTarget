package middlewares

import (
	"net/http"
	"testing"
)

func TestIsAuthorized(t *testing.T) {
	t.Parallel()
	t.Fatal("not implemeted")

	type testCase struct {
		desc string
		stub http.HandlerFunc
	}

	testCases := []testCase{{
		desc: "should add nil when no cookie",
		stub: nil,
	}, {
		desc: "should add nil when there's cookie in session storage",
		stub: nil,
	}, {
		desc: "should add user object when cookie is good",
		stub: nil,
	}}

	for _, tc := range testCases {

	}
}
