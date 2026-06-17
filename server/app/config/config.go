package config

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/cornbuddy/reflectiveTarget/server/infra/daos"
	"github.com/cornbuddy/reflectiveTarget/server/infra/log"
	"github.com/cornbuddy/reflectiveTarget/server/infra/repositories"
	"github.com/cornbuddy/reflectiveTarget/server/infra/utils"
)

type Config struct {
	daos.HealthDao
	daos.SessionStore
	daos.ShotsDao
	daos.UserDao
	repositories.TargetRepo
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
	if err := retry(ctx, "connect to cache", func() error {
		cache = redis.NewClient(&redis.Options{
			Addr: cfg.CacheAddr,
		})

		return cache.Ping(ctx).Err()
	}); err != nil {
		return nil, err
	}

	var db *sql.DB
	dbConnStr := fmt.Sprintf(
		"postgresql://%s:%s@%s:5432/%s?sslmode=disable",
		cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.Db,
	)
	if err := retry(ctx, "connect to db", func() error {
		var err error
		db, err = sql.Open("pgx", dbConnStr)
		if err != nil {
			return err
		}

		return db.Ping()
	}); err != nil {
		return nil, err
	}

	if err := utils.InitDatabase(db); err != nil {
		return nil, err
	}

	return &Config{
		HealthDao:    daos.HealthDao{DB: db, Cache: cache},
		SessionStore: daos.SessionStore{Cache: cache},
		ShotsDao:     daos.ShotsDao{DB: db},
		UserDao:      daos.UserDao{DB: db},
		TargetRepo:   repositories.TargetRepo{DB: db},
	}, nil
}

func retry(ctx context.Context, operation string, f func() error) error {
	const attempts = 5
	const delay = 3 * time.Second

	var err error
	for attempt := range attempts {
		log.Info(ctx, operation, zap.Int("attempt", attempt))
		if err = f(); err != nil {
			log.Warn(ctx, operation,
				zap.String("status", "failed"),
				zap.Duration("delay", delay),
			)
			time.Sleep(delay)
		} else {
			log.Info(ctx, operation, zap.String("status", "success"))
			break
		}
	}

	return err
}
