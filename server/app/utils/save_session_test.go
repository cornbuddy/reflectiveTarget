package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveSessionShouldUpdateSessionStore(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	wantAuth := false
	token, err := SaveSession(ctx, store, wantAuth, w)
	require.NoError(t, err)

	header := w.Header().Get("Set-Cookie")
	cookie, err := http.ParseSetCookie(header)
	require.NoError(t, err)
	require.NotNil(t, cookie)
	assert.Equal(t, token, cookie.Value)

	auth, err := store.IsAuthenitcated(ctx, token)
	require.NoError(t, err)
	assert.NotNil(t, auth)
	assert.Equal(t, wantAuth, *auth, "should update session store")
}
