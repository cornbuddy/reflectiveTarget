package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogMiddlewareShouldNotChangeRequest(t *testing.T) {
	t.Parallel()

	want := httptest.NewRequest(http.MethodGet, "/", nil)
	stub := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, want, r)
	})

	mw.Log(stub).ServeHTTP(httptest.NewRecorder(), want)
}
