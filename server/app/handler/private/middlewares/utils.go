package middlewares

import (
	"net/http"

	"go.uber.org/zap"

	. "github.com/cornbuddy/reflectiveTarget/server/infra/logger"
)

func internalServerError(w http.ResponseWriter, msg string) {
	Log.Error("unexpected error in middleware chain",
		zap.String("error", msg),
	)
	http.Error(w, msg, http.StatusInternalServerError)
}
