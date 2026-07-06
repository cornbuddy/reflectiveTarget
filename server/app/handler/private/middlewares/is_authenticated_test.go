package middlewares_test

import (
	"net/http"
	"testing"

	"github.com/bloomberg/go-testgroup"
	"github.com/gorilla/mux"

	"github.com/cornbuddy/reflectiveTarget/server/app/constants"
	"github.com/cornbuddy/reflectiveTarget/server/app/session"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

type IsAuthenticatedSuite struct {
	sessionID session.SessionID
	handler   http.HandlerFunc
}

func (s *IsAuthenticatedSuite) ShouldErrorIfNoToken(t *testgroup.T) {
	r, _, err := utils.MakeRequest("", http.MethodGet, "/", s.handler, nil)
	t.Require.NoError(err)
	t.Equal(http.StatusForbidden, r.StatusCode)
}

func (s *IsAuthenticatedSuite) ShouldErrorIfTokenIsInvalid(t *testgroup.T) {
	cookies := []*http.Cookie{{
		Name:  constants.SessionCookieName,
		Value: "kek",
	}}
	r, _, err := utils.MakeRequestWithCookies(
		"", http.MethodGet, "/", s.handler, nil, cookies...,
	)
	t.Require.NoError(err)
	t.Equal(http.StatusForbidden, r.StatusCode)
}

func (s *IsAuthenticatedSuite) Should200IfAuthenticated(t *testgroup.T) {
	cookies := []*http.Cookie{{
		Name:  constants.SessionCookieName,
		Value: s.sessionID.String(),
	}}
	r, _, err := utils.MakeRequestWithCookies(
		"", http.MethodGet, "/", s.handler, nil, cookies...,
	)
	t.Require.NoError(err)
	t.Equal(http.StatusOK, r.StatusCode)
}

func TestIsAuthenticated(t *testing.T) {
	t.Parallel()

	testgroup.RunInParallel(t, new(IsAuthenticatedSuite))
}

func (s *IsAuthenticatedSuite) PreGroup(t *testgroup.T) {
	id := session.MakeID()
	data := session.Data{
		UserID:   69,
		Username: username,
	}
	t.Require.NoError(store.Update(ctx, id, data))

	s.sessionID = id
	s.handler = chain(
		emptyStub,
		mw.IsAuthenticated,
		mw.SaveSession,
	).ServeHTTP
}

func chain(mux http.Handler, mwf ...mux.MiddlewareFunc) http.Handler {
	for _, mw := range mwf {
		mux = mw(mux)
	}

	return mux
}
