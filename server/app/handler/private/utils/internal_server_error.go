package utils

import (
	"net/http"

	"go.uber.org/zap"
)

func InternalServerError(
	log *zap.Logger, w http.ResponseWriter, msg string, err error,
) {
	log.Error(msg, zap.Error(err))
	http.Error(w, msg, http.StatusInternalServerError)
}
