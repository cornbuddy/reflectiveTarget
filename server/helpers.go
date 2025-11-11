package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/cornbuddy/reflectiveTarget/server/private/daos"
	"github.com/cornbuddy/reflectiveTarget/server/private/handlers"
	"github.com/cornbuddy/reflectiveTarget/server/private/routers"
	"github.com/cornbuddy/reflectiveTarget/server/private/utils"
)

type Config struct {
	*sql.DB
	daos.UserDao
}

var ErrNoEnvVar = fmt.Errorf("no environment variable")

func MakeMux(config *Config) http.Handler {
	health := handlers.HealthHandler{
		DB: config.DB,
	}
	authz := handlers.AuthzHandler{
		UserDao: config.UserDao,
	}
	index := handlers.IndexHandler{}

	views := routers.ViewsRouter{
		AuthzHandler: authz,
		IndexHandler: index,
	}

	mux := http.NewServeMux()
	mux.Handle("/", views.Routes())
	mux.HandleFunc("GET /api/health", health.Get)

	return mux
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

	const retries = 5
	const delay = 3 * time.Second
	var db *sql.DB
	var err error

	for attempt := range retries {
		log.Printf("connecting to db, attempt #%d", attempt)
		db, err = sql.Open("pgx", connStr)
		connected := db.Ping() == nil
		if connected {
			log.Println("connected to db")
			break
		}

		log.Printf("failed, waiting %v...", delay)
		time.Sleep(delay)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
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
