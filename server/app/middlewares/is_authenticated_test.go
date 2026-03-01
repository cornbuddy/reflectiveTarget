package middlewares

import (
	"net/http"
	"testing"

	"github.com/bloomberg/go-testgroup"

	"github.com/cornbuddy/reflectiveTarget/server/app/constants"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

func (s *IsAuthenticatedSuite) ShouldErrorIfNoToken(t *testgroup.T) {
	r := utils.MakeRequest("", http.MethodGet, "/", s.handler, nil)
	t.Equal(http.StatusForbidden, r.StatusCode)
}

func (s *IsAuthenticatedSuite) ShouldErrorIfTokenIsInvalid(t *testgroup.T) {
	cookies := []*http.Cookie{{
		Name:  constants.SessionCookieName,
		Value: "kek",
	}}
	r := utils.MakeRequestWithCookies(
		"", http.MethodGet, "/", s.handler, nil, cookies...,
	)
	t.Equal(http.StatusForbidden, r.StatusCode)
}

func (s *IsAuthenticatedSuite) Should200IfAuthenticated(t *testgroup.T) {
	cookies := []*http.Cookie{{
		Name:  constants.SessionCookieName,
		Value: s.authenticatedToken,
	}}
	r := utils.MakeRequestWithCookies(
		"", http.MethodGet, "/", s.handler, nil, cookies...,
	)
	t.Equal(http.StatusOK, r.StatusCode)
}

func TestIsAuthenticated(t *testing.T) {
	t.Parallel()

	testgroup.RunInParallel(t, new(IsAuthenticatedSuite))
}
