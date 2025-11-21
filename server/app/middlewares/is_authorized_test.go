package middlewares

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/app/handlers"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

func TestIsAuthorizedShouldAddToCtxIfCookieIsInTheSessionStore(t *testing.T) {
	t.Parallel()

	cookie := "cookie"
	username := "username"
	require.NoError(t, store.SaveSessionForUser(username, cookie))

	stub := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got := r.Context().Value(Username)
		require.NotNil(t, got)
		assert.Equal(t, username, got.(string))
	})
	cookies := []*http.Cookie{{
		Name:  handlers.SessionCookieName,
		Value: cookie,
	}}
	handler := mv.IsAuthorized(stub).ServeHTTP
	utils.MakeRequestWithCookies(
		"", http.MethodGet, "/", handler, nil, cookies...,
	)
}

func TestIsAuthorizedShouldAddNilToCtxIfNoCookieInRequest(t *testing.T) {
	t.Parallel()

	stub := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Nil(t, r.Context().Value(Username))
	})

	handler := mv.IsAuthorized(stub).ServeHTTP
	utils.MakeRequest("", http.MethodGet, "/", handler, nil)
}
