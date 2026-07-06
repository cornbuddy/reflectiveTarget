package middlewares_test

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/app/constants"
	"github.com/cornbuddy/reflectiveTarget/server/app/session"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

func TestSaveSessionShouldPutSessionDataToCtx(t *testing.T) {
	t.Parallel()

	type testCase struct {
		desc      string
		sessionId session.SessionID
		want      session.Data
	}

	authorizedID := session.MakeID()
	want := session.Data{
		UserID:   69,
		Username: username,
	}
	require.NoError(t, store.Update(ctx, authorizedID, want))

	testCases := []testCase{{
		"should put cached data to context if any",
		authorizedID,
		want,
	}, {
		"should put zero object on cache miss",
		session.MakeID(),
		session.Data{},
	}}

	for _, tc := range testCases {
		cookies := []*http.Cookie{{
			Name:  constants.SessionCookieName,
			Value: tc.sessionId.String(),
		}}
		stub := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.NotNil(t, w)
			require.NotNil(t, r)

			got := session.Read(r.Context())
			assert.EqualExportedValues(t, tc.want, got, tc.desc)
		})
		handler := mw.SaveSession(stub).ServeHTTP
		_, _, err := utils.MakeRequestWithCookies(
			"", http.MethodGet, "/", handler, nil, cookies...,
		)
		require.NoError(t, err)
	}
}

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
	_, _, err := utils.MakeRequestWithCookies(
		"", http.MethodGet, "/", handler, nil, cookies...,
	)
	require.NoError(t, err)
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
	_, _, err := utils.MakeRequest("", http.MethodGet, "/", handler, nil)
	require.NoError(t, err)
}

func TestSaveSessionShouldRespectExistingSessionToken(t *testing.T) {
	t.Parallel()

	id := session.MakeID()
	require.NoError(t, store.Update(ctx, id, session.Data{}))

	stub := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := w.Header().Get("Set-Cookie")
		assert.Empty(t, header)
	})

	handler := mw.SaveSession(stub).ServeHTTP
	cookies := []*http.Cookie{{
		Name:  constants.SessionCookieName,
		Value: id.String(),
	}}
	_, _, err := utils.MakeRequestWithCookies(
		"", http.MethodGet, "/", handler, nil, cookies...,
	)
	require.NoError(t, err)
}
