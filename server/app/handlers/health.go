package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type HealthResponse struct {
	DbConnected      bool `json:"database-connected"`
	DbConnections    int  `json:"database-connections"`
	CacheConnected   bool `json:"cache-connected"`
	CacheConnections int  `json:"cache-connections"`
}

type healthHandler struct {
	*sql.DB
	Cache *redis.Client
}

func (h healthHandler) get(resp http.ResponseWriter, req *http.Request) {
	dbConnections := 0
	dbConnected := h.DB.Ping() == nil
	if dbConnected {
		dbConnections = h.DB.Stats().OpenConnections
	}

	cacheConnections := 0
	cacheConnected := h.Cache.Ping(req.Context()).Err() == nil
	if cacheConnected {
		cacheConnections = int(h.Cache.PoolStats().TotalConns)
	}

	var status int
	if dbConnected && cacheConnected {
		status = http.StatusOK
	} else {
		status = http.StatusServiceUnavailable
	}

	hr, _ := json.Marshal(HealthResponse{
		DbConnected:      dbConnected,
		DbConnections:    dbConnections,
		CacheConnected:   cacheConnected,
		CacheConnections: cacheConnections,
	})
	resp.WriteHeader(status)
	resp.Write(hr)
}
