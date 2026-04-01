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

func TestSaveSessionShouldPutSessionDataToCtx(t *testing.T) {
	t.Parallel()

	type testCase struct {
		desc      string
		sessionId string
		want      sessiondata.SessionData
	}

	sessionId := "kekeke"
	want := sessiondata.SessionData{
		IsAuthenticated: true,
		UserID:          69,
		Username:        "kek",
	}
	require.NoError(t, store.Update(ctx, sessionId, want))

	testCases := []testCase{{
		"should put cached data to context if any",
		sessionId,
		want,
	}, {
		"should put zero object on cache miss",
		"not found",
		sessiondata.SessionData{},
	}}

	for _, tc := range testCases {
		cookies := []*http.Cookie{{
			Name:  constants.SessionCookieName,
			Value: tc.sessionId,
		}}
		stub := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.NotNil(t, w)
			require.NotNil(t, r)

			got := sessiondata.Read(r.Context())
			assert.EqualExportedValues(t, tc.want, got, tc.desc)
		})
		handler := mw.SaveSession(stub).ServeHTTP
		utils.MakeRequestWithCookies(
			"", http.MethodGet, "/", handler, nil, cookies...,
		)
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
