package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

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
		log.Fatal("failed to run db: %w", err)
	}

	defer func() {
		if err := dbCont.Terminate(ctx); err != nil {
			log.Fatal("failed to stop db: %w", err)
		}
	}()

	dbHost, err = dbCont.ContainerIP(ctx)
	if err != nil {
		log.Fatal("failed to fetch db ip: %w", err)
	}

	cacheCont, err := tcredis.Run(ctx, utils.CacheImage)
	if err != nil {
		log.Fatal("failed to run cache: %w", err)
	}

	defer func() {
		if err := cacheCont.Terminate(ctx); err != nil {
			log.Fatal("failed to stop cache: %w", err)
		}
	}()

	cacheIp, err := cacheCont.ContainerIP(ctx)
	if err != nil {
		log.Fatal("failed to fetch cache ip: %w", err)
	}

	cacheAddr = fmt.Sprintf("%s:6789", cacheIp)

	os.Exit(m.Run())
}
