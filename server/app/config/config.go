package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/caarlos0/env/v11"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"

	"github.com/cornbuddy/reflectiveTarget/server/infra/daos"
	"github.com/cornbuddy/reflectiveTarget/server/infra/utils"
)

type Config struct {
	daos.HealthDao
	daos.SessionStore
	daos.ShotsDao
	daos.UserDao
}

func MakeConfig(ctx context.Context) (*Config, error) {
	type config struct {
		CacheAddr  string `env:"CACHE_ADDRESS,notEmpty,required"`
		DbPassword string `env:"PGPASSWORD,notEmpty,required"`
		DbUser     string `env:"PGUSER,notEmpty,required"`
		Db         string `env:"PGDATABASE,notEmpty,required"`
		DbHost     string `env:"PGHOST,notEmpty,required"`
	}

	var cfg config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	var cache *redis.Client
	retry("connect to cache", func() error {
		cache = redis.NewClient(&redis.Options{
			Addr: cfg.CacheAddr,
		})

		if err := cache.Ping(ctx).Err(); err != nil {
			return err
		}

		return nil
	})

	dbConnStr := fmt.Sprintf(
		"postgresql://%s:%s@%s:5432/%s?sslmode=disable",
		cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.Db,
	)

	var db *sql.DB
	retry("connect to db", func() error {
		var err error
		db, err = sql.Open("pgx", dbConnStr)
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
		HealthDao:    daos.HealthDao{Ctx: ctx, DB: db, Cache: cache},
		SessionStore: daos.SessionStore{Ctx: ctx, Cache: cache},
		ShotsDao:     daos.ShotsDao{DB: db},
		UserDao:      daos.UserDao{DB: db},
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
