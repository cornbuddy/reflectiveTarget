package config_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

var (
	dbHost    string
	cacheAddr string

	ctx      = context.TODO()
	username = "test_user"
	password = "kekekeke"
	database = "testdb"
)

func TestMain(m *testing.M) {
	str := wait.ForLog("database system is ready to accept connections").
		WithOccurrence(2).
		WithStartupTimeout(5 * time.Second)
	opts := []tc.ContainerCustomizer{
		tc.WithWaitStrategy(str),
		postgres.WithUsername(username),
		postgres.WithPassword(password),
		postgres.WithDatabase(database),
	}

	dbCont, err := postgres.Run(ctx, utils.DbImage, opts...)
	if err != nil {
		log.Fatalf("failed to run db: %v", err)
	}

	defer func() {
		if err := dbCont.Terminate(ctx); err != nil {
			log.Fatalf("failed to stop db: %v", err)
		}
	}()

	dbHost, err = dbCont.ContainerIP(ctx)
	if err != nil {
		log.Fatalf("failed to fetch db ip: %v", err)
	}

	cacheCont, err := tcredis.Run(ctx, utils.CacheImage)
	if err != nil {
		log.Fatalf("failed to run cache: %v", err)
	}

	defer func() {
		if err := cacheCont.Terminate(ctx); err != nil {
			log.Fatalf("failed to stop cache: %v", err)
		}
	}()

	uri, err := cacheCont.ConnectionString(ctx)
	if err != nil {
		log.Fatalf("failed to fetch uri for cache: %v", err)
	}

	cacheOpts, err := redis.ParseURL(uri)
	if err != nil {
		log.Fatalf("failed to parse cache uri: %v", err)
	}

	cacheAddr = cacheOpts.Addr

	os.Exit(m.Run())
}
