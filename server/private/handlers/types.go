package handlers

import (
	"database/sql"

	"github.com/cornbuddy/reflectiveTarget/server/private/daos"
)

type AuthzRouter struct {
	daos.UserDao
}

type HealthResponse struct {
	Connected   bool `json:"connected"`
	Connections int  `json:"connections"`
}

type HealthRouter struct {
	*sql.DB
}
