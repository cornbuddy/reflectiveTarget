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

	var db *sql.DB
	retry("connect to db", func() error {
		var err error
		db, err = sql.Open("pgx", connStr)
		if err != nil {
			return err
		}

		if err := db.Ping(); err != nil {
			return err
		}

		return nil
	})

	if err := utils.InitDatabase(db); err != nil {
		return nil, err
	}

	return &Config{
		DB:       db,
		UserDao:  daos.UserDao{DB: db},
		ShotsDao: daos.ShotsDao{DB: db},
	}, nil
}

const attempts = 5
const delay = 3 * time.Second

func retry(operation string, f func() error) error {
	var err error
	for attempt := range attempts {
		log.Printf("%s, attempt #%d", operation, attempt)
		if err = f(); err != nil {
			log.Printf("%s failed, waiting %v...", operation, delay)
			time.Sleep(delay)
		} else {
			log.Printf("%s is succeeded", operation)
			break
		}
	}

	return err
}
