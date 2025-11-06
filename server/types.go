package main

import (
	"database/sql"

	"github.com/cornbuddy/reflectiveTarget/server/private/daos"
)

type Config struct {
	sql.DB
	daos.UserDao
}
