package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type HealthResponse struct {
	Connected   bool `json:"connected"`
	Connections int  `json:"connections"`
}

type HealthRouter struct {
	*sql.DB
}

func (r HealthRouter) Get(resp http.ResponseWriter, req *http.Request) {
	var connections, status int
	connected := r.DB.Ping() == nil
	if connected {
		connections = r.DB.Stats().OpenConnections
		status = http.StatusOK
	} else {
		connections = 0
		status = http.StatusServiceUnavailable
	}

	hr, _ := json.Marshal(HealthResponse{
		Connected:   connected,
		Connections: connections,
	})
	resp.WriteHeader(status)
	resp.Write(hr)
}
