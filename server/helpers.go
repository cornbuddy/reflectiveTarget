package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/cornbuddy/reflectiveTarget/server/private/daos"
	"github.com/cornbuddy/reflectiveTarget/server/private/utils"
)

type Config struct {
	*sql.DB
	daos.UserDao
}

var ErrNoEnvVar = fmt.Errorf("no environment variable")

func MakeMux() http.Handler {
	return nil
}

func MakeConfig() (*Config, error) {
	var password, user, database, host string

	if pwd, ok := os.LookupEnv("PGPASSWORD"); !ok {
		return nil, fmt.Errorf("%w: %s", ErrNoEnvVar, "PGPASSWORD")
	} else {
		password = pwd
	}

	if usr, ok := os.LookupEnv("PGUSER"); !ok {
		return nil, fmt.Errorf("%w: %s", ErrNoEnvVar, "PGUSER")
	} else {
		user = usr
	}

	if db, ok := os.LookupEnv("PGDATABASE"); !ok {
		return nil, fmt.Errorf("%w: %s", ErrNoEnvVar, "PGDATABASE")
	} else {
		database = db
	}

	if h, ok := os.LookupEnv("PGHOST"); !ok {
		return nil, fmt.Errorf("%w: %s", ErrNoEnvVar, "PGHOST")
	} else {
		host = h
	}

	connStr := fmt.Sprintf(
		"postgresql://%s:%s@%s:5432/%s?sslmode=disable",
		user, password, host, database,
	)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	if err := utils.InitDatabase(db); err != nil {
		return nil, err
	}

	return &Config{
		DB: db,
		UserDao: daos.UserDao{
			DB: db,
		},
	}, nil
}
