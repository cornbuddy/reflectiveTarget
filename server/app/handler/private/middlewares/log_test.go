package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cornbuddy/reflectiveTarget/server/app/constants"
)

func TestLoggerMiddlewareShouldAddLoggerToContext(t *testing.T) {
	t.Parallel()

	stub := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Context().Value(constants.RequestIDCtx)
		assert.NotEmpty(t, reqID)
		assert.IsType(t, reqID, "string")
	})

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	mw.Logger(stub).ServeHTTP(httptest.NewRecorder(), r)
}
