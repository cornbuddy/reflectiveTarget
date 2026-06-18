package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/cornbuddy/reflectiveTarget/server/infra/log"
)

func TestLoggerMiddlewareShouldAddLoggerToContext(t *testing.T) {
	t.Parallel()

	stub := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger, ok := r.Context().Value(log.LoggerCtx).(*zap.Logger)
		assert.True(t, ok)
		assert.NotEmpty(t, logger)
	})

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	mw.Logger(stub).ServeHTTP(httptest.NewRecorder(), r)
}
