package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsAuthenticatedShould403IfNotAuthenticated(t *testing.T) {
	t.Parallel()

	stub := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	mw.IsAuthenticated(stub).ServeHTTP(w, r)

	assert.Equal(t, http.StatusForbidden, w.Result().StatusCode)
}
