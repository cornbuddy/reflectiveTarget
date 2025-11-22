package middlewares

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/app/handlers"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

func TestSaveSessionShouldRespectExistingSessionToken(t *testing.T) {
	t.Parallel()

	token := "i'm a session token"
	require.NoError(t, store.SaveSession(token, false))

	stub := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(handlers.SessionCookieName)
		require.NoError(t, err)
		require.NotNil(t, cookie)
		assert.Equal(t, token, cookie.Value)
	})

	handler := mw.SaveSession(stub).ServeHTTP
	cookies := []*http.Cookie{{
		Name:  handlers.SessionCookieName,
		Value: "kek",
	}}
	utils.MakeRequestWithCookies(
		"", http.MethodGet, "/", handler, nil, cookies...,
	)
}

func TestSaveSessionShouldAddSessionCookieIfNotPresent(t *testing.T) {
	t.Parallel()

	stub := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(handlers.SessionCookieName)
		require.NoError(t, err)
		require.NotNil(t, cookie)
		assert.NotEmpty(t, cookie.Value)
	})

	handler := mw.SaveSession(stub).ServeHTTP
	utils.MakeRequest("", http.MethodGet, "/", handler, nil)
}
