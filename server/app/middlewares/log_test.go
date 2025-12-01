package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cornbuddy/reflectiveTarget/server/app/constants"
	. "github.com/cornbuddy/reflectiveTarget/server/infra/logger"
)

func TestLoggerMiddlewareShouldAddLoggerToContext(t *testing.T) {
	t.Parallel()

	stub := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := r.Context().Value(constants.LoggerCtx)
		assert.NotNil(t, log)
		assert.IsType(t, Log, log)
	})

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	mw.Logger(stub).ServeHTTP(httptest.NewRecorder(), r)
}
