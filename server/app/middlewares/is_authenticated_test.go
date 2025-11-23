package middlewares

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/app/constants"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

func TestIsAuthorizedShouldAddToCtxIfCookieIsInTheSessionStore(t *testing.T) {
	t.Parallel()

	token := "cookie"
	authenticated := false
	require.NoError(t, store.SaveSession(token, authenticated))

	stub := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got := r.Context().Value(constants.AuthenticatedCtx).(*bool)
		require.NotNil(t, got)
		assert.Equal(t, authenticated, *got)
	})
	cookies := []*http.Cookie{{
		Name:  constants.SessionCookieName,
		Value: token,
	}}
	handler := mw.IsAuthenticated(stub).ServeHTTP
	utils.MakeRequestWithCookies(
		"", http.MethodGet, "/", handler, nil, cookies...,
	)
}

func TestIsAuthorizedShouldAddNilToCtxIfNoCookieInRequest(t *testing.T) {
	t.Parallel()

	stub := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Nil(t, r.Context().Value(constants.AuthenticatedCtx))
	})

	handler := mw.IsAuthenticated(stub).ServeHTTP
	utils.MakeRequest("", http.MethodGet, "/", handler, nil)
}
