package utils

import (
	"net/http"

	"go.uber.org/zap"
)

func InternalServerError(
	log *zap.Logger, w http.ResponseWriter, msg string, err error,
) {
	httpError(log, w, http.StatusInternalServerError, msg, err)
}

func BadRequest(
	log *zap.Logger, w http.ResponseWriter, msg string, err error,
) {
	httpError(log, w, http.StatusBadRequest, msg, err)
}

func NotFound(
	log *zap.Logger, w http.ResponseWriter, msg string, err error,
) {
	httpError(log, w, http.StatusNotFound, msg, err)
}

func httpError(
	log *zap.Logger, w http.ResponseWriter, code int, msg string, err error,
) {
	log.Error(msg, zap.Error(err))
	http.Error(w, msg, code)
}
