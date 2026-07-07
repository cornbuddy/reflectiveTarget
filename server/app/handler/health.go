package handler

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/utils"
	"github.com/cornbuddy/reflectiveTarget/server/infra/daos"
	"github.com/cornbuddy/reflectiveTarget/server/infra/log"
)

type HealthResponse struct {
	daos.HealthStatus `json:",inline"`
}

type healthHandler struct {
	dao daos.HealthDao
}

func (h healthHandler) get(w http.ResponseWriter, r *http.Request) {
	status := h.dao.CheckHealth(r.Context())
	code := http.StatusInternalServerError
	if status.CacheConnected && status.DbConnected {
		code = http.StatusOK
	}

	log := log.Logger(r.Context())
	hr, err := json.Marshal(HealthResponse{
		HealthStatus: status,
	})
	if err != nil {
		log.Error("failed to marshal health status", zap.Error(err))
		utils.HttpError(w, http.StatusInternalServerError)

		return
	}

	w.WriteHeader(code)
	utils.LogBadWrites(log)(w.Write(hr))
}
