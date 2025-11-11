package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/cornbuddy/reflectiveTarget/server/private/model"
)

func (r HealthHandler) Get(resp http.ResponseWriter, req *http.Request) {
	var connections, status int
	connected := r.DB.Ping() == nil
	if connected {
		connections = r.DB.Stats().OpenConnections
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
