package utils_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/utils"
	"github.com/cornbuddy/reflectiveTarget/server/app/session"
)

func TestSaveSessionShouldUpdateSessionStore(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	wantSession := session.Data{Username: "kek"}
	token, err := utils.SaveSession(ctx, store, wantSession, w)
	require.NoError(t, err)

	header := w.Header().Get("Set-Cookie")
	cookie, err := http.ParseSetCookie(header)
	require.NoError(t, err)
	require.NotNil(t, cookie)
	assert.Equal(t, token, cookie.Value)
	assert.Equal(t, http.SameSiteStrictMode, cookie.SameSite)
	assert.True(t, cookie.HttpOnly)

	gotSession, err := store.Get(ctx, token)
	require.NoError(t, err)
	assert.NotNil(t, gotSession)
	assert.Equal(t, wantSession, *gotSession, "should update session store")
}
