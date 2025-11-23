package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/caarlos0/env/v11"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/cornbuddy/reflectiveTarget/server/app/handlers"
	"github.com/cornbuddy/reflectiveTarget/server/infra/daos"
	"github.com/cornbuddy/reflectiveTarget/server/infra/utils"
)

type Config struct {
	*sql.DB
	daos.UserDao
	daos.ShotsDao
}

func MakeMux(config *Config) http.Handler {
	views := handlers.ViewsRouter{
		UserDao: config.UserDao,
	}
	api := handlers.ApiRouter{
		DB:       config.DB,
		ShotsDao: config.ShotsDao,
	}

	mux := http.NewServeMux()
	mux.Handle("/", views.Routes())
	mux.Handle("/api/", http.StripPrefix("/api", api.Routes()))

	return mux
}

func MakeConfig() (*Config, error) {
	type config struct {
		DbPassword string `env:"PGPASSWORD,notEmpty,required"`
		DbUser     string `env:"PGUSER,notEmpty,required"`
		Db         string `env:"PGDATABASE,notEmpty,required"`
		DbHost     string `env:"PGHOST,notEmpty,required"`
	}

	var cfg config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	connStr := fmt.Sprintf(
		"postgresql://%s:%s@%s:5432/%s?sslmode=disable",
		cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.Db,
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
		DB:       db,
		UserDao:  daos.UserDao{DB: db},
		ShotsDao: daos.ShotsDao{DB: db},
	}, nil
}
