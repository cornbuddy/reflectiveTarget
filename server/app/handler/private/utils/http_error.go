package utils

import (
	"net/http"

	"go.uber.org/zap"
)

type logger func(string, ...zap.Field)

func InternalServerError(
	log *zap.Logger, w http.ResponseWriter, msg string, err error,
) {
	httpError(log.Error, w, http.StatusInternalServerError, msg, err)
}

func BadRequest(
	log *zap.Logger, w http.ResponseWriter, msg string, err error,
) {
	httpError(log.Warn, w, http.StatusBadRequest, msg, err)
}

func NotFound(
	log *zap.Logger, w http.ResponseWriter, msg string, err error,
) {
	httpError(log.Warn, w, http.StatusNotFound, msg, err)
}

func httpError(
	log logger, w http.ResponseWriter, code int, msg string, err error,
) {
	log(msg, zap.Error(err))
	http.Error(w, msg, code)
}
