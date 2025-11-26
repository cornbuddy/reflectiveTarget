package daos

import (
	"context"
	"database/sql"

	"github.com/redis/go-redis/v9"
)

type HealthDao struct {
	Ctx context.Context
	*sql.DB
	Cache *redis.Client
}

type HealthStatus struct {
	DbConnected      bool `json:"database-connected"`
	DbConnections    int  `json:"database-connections"`
	CacheConnected   bool `json:"cache-connected"`
	CacheConnections int  `json:"cache-connections"`
}

func (dao HealthDao) CheckHealth() HealthStatus {
	dbConnections := 0
	dbConnected := dao.DB.Ping() == nil
	if dbConnected {
		dbConnections = dao.DB.Stats().OpenConnections
	}

	cacheConnections := 0
	cacheConnected := dao.Cache.Ping(dao.Ctx).Err() == nil
	if cacheConnected {
		cacheConnections = int(dao.Cache.PoolStats().TotalConns)
	}

	return HealthStatus{
		DbConnected:      dbConnected,
		DbConnections:    dbConnections,
		CacheConnected:   cacheConnected,
		CacheConnections: cacheConnections,
	}
}
