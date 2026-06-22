package config

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
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

const (
	DefaultTimeout         = 15 * time.Second
	DefaultPort            = 8080
	ReadTimeout            = 5 * time.Second
	WriteTimeoutMultiplier = 2
	IdleTiemoutMultiplier  = 3

	minPort = 1024
	maxPort = 65535
)

type Config struct {
	daos.HealthDao
	daos.SessionStore
	daos.ShotsDao
	daos.UserDao
	repositories.TargetRepo

	Timeout time.Duration
	Port    int
}

func MakeConfig(ctx context.Context) (*Config, error) {
	type config struct {
		TimeoutSecs string `env:"TIMEOUT_SECONDS" envDefault:"15"`
		Port        string `env:"PORT" envDefault:"8080"`
		CacheAddr   string `env:"CACHE_ADDRESS,notEmpty,required"`
		DbPassword  string `env:"PGPASSWORD,notEmpty,required"`
		DbUser      string `env:"PGUSER,notEmpty,required"`
		Db          string `env:"PGDATABASE,notEmpty,required"`
		DbHost      string `env:"PGHOST,notEmpty,required"`
	}

	log := log.Logger(ctx)

	var cfg config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	var timeout time.Duration
	secs, err := strconv.ParseInt(cfg.TimeoutSecs, 10, 64)
	switch {
	case err != nil:
		log.Warn(
			"failed to parse timeout environment variable",
			zap.String("value", cfg.TimeoutSecs), zap.Error(err),
		)
		timeout = DefaultTimeout
	case secs <= 0:
		log.Warn("timeout value should be positive", zap.Int64("value", secs))
		timeout = DefaultTimeout
	default:
		timeout = time.Duration(secs) * time.Second
	}

	var port int
	p, err := strconv.Atoi(cfg.Port)
	switch {
	case err != nil:
		log.Warn(
			"failed to parse port environment variable",
			zap.String("value", cfg.Port), zap.Error(err),
		)
		port = DefaultPort
	case p < minPort || p > maxPort:
		log.Warn("bad port", zap.Int("port", p))
		port = DefaultPort
	default:
		port = p
	}

	var cache *redis.Client
	if err := retry(log, "connect to cache", func() error {
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
	if err := retry(log, "connect to db", func() error {
		var err error
		db, err = sql.Open("pgx", dbConnStr)
		if err != nil {
			return err
		}

		return db.PingContext(ctx)
	}); err != nil {
		return nil, err
	}

	if err := utils.InitDatabase(db); err != nil {
		return nil, err
	}

	log.Info("config is created")

	return &Config{
		Timeout:      timeout,
		Port:         port,
		HealthDao:    daos.HealthDao{DB: db, Cache: cache},
		SessionStore: daos.SessionStore{Cache: cache},
		ShotsDao:     daos.ShotsDao{DB: db},
		UserDao:      daos.UserDao{DB: db},
		TargetRepo:   repositories.TargetRepo{DB: db},
	}, nil
}

func retry(log *zap.Logger, operation string, f func() error) error {
	const attempts = 5
	const delay = 3 * time.Second

	var err error
	for attempt := range attempts {
		log := log.With(zap.Int("attempt", attempt))
		if err = f(); err != nil {
			log.Warn(operation,
				zap.String("status", "failed"),
				zap.Duration("delay", delay),
			)
			time.Sleep(delay)
		} else {
			log.Info(operation, zap.String("status", "success"))

			break
		}
	}

	return err
}
