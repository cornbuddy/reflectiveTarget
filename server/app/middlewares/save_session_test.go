package middlewares

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/app/constants"
	"github.com/cornbuddy/reflectiveTarget/server/app/sessiondata"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

func TestSaveSessionShouldResetRequestCookieWhenItsNotInStore(t *testing.T) {
	t.Parallel()

	token := uuid.NewString()
	stub := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := w.Header().Get("Set-Cookie")
		cookie, err := http.ParseSetCookie(header)
		require.NoError(t, err)
		require.NotNil(t, cookie)
		assert.NoError(t, uuid.Validate(cookie.Value))
		assert.NotEqual(t, token, cookie.Value)
	})

	handler := mw.SaveSession(stub).ServeHTTP
	cookies := []*http.Cookie{{
		Name:  constants.SessionCookieName,
		Value: token,
	}}
	utils.MakeRequestWithCookies(
		"", http.MethodGet, "/", handler, nil, cookies...,
	)
}

func TestSaveSessionShouldAddSessionCookieIfNotPresent(t *testing.T) {
	t.Parallel()

	stub := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := w.Header().Get("Set-Cookie")
		cookie, err := http.ParseSetCookie(header)
		require.NoError(t, err)
		require.NotNil(t, cookie)
		assert.NoError(t, uuid.Validate(cookie.Value))
	})

	handler := mw.SaveSession(stub).ServeHTTP
	utils.MakeRequest("", http.MethodGet, "/", handler, nil)
}

func TestSaveSessionShouldRespectExistingSessionToken(t *testing.T) {
	t.Parallel()

	token := uuid.NewString()
	require.NoError(t, store.Update(ctx, token, sessiondata.SessionData{}))

	stub := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := w.Header().Get("Set-Cookie")
		assert.Empty(t, header)
	})

	handler := mw.SaveSession(stub).ServeHTTP
	cookies := []*http.Cookie{{
		Name:  constants.SessionCookieName,
		Value: token,
	}}
	utils.MakeRequestWithCookies(
		"", http.MethodGet, "/", handler, nil, cookies...,
	)
}
