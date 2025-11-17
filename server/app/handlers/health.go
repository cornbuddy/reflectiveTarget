package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/cornbuddy/reflectiveTarget/server/private/model"
)

type healthHandler struct {
	*sql.DB
}

func (h healthHandler) get(resp http.ResponseWriter, req *http.Request) {
	var connections, status int
	connected := h.DB.Ping() == nil
	if connected {
		connections = h.DB.Stats().OpenConnections
		status = http.StatusOK
	} else {
		connections = 0
		status = http.StatusServiceUnavailable
	}

	hr, _ := json.Marshal(model.HealthResponse{
		Connected:   connected,
		Connections: connections,
	})
	resp.WriteHeader(status)
	resp.Write(hr)
}
