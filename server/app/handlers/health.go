package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/cornbuddy/reflectiveTarget/server/infra/daos"
)

type HealthResponse struct {
	daos.HealthStatus `json:",inline"`
}

type healthHandler struct {
	daos.HealthDao
}

func (h healthHandler) get(resp http.ResponseWriter, req *http.Request) {
	status := h.HealthDao.CheckHealth()
	code := http.StatusInternalServerError
	if status.CacheConnected && status.DbConnected {
		code = http.StatusOK
	}

	hr, _ := json.Marshal(HealthResponse{
		HealthStatus: status,
	})
	resp.WriteHeader(code)
	resp.Write(hr)
}
