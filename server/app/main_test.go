package main

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

var (
	host string

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

	cont, err := postgres.Run(ctx, utils.DbImage, opts...)
	if err != nil {
		log.Fatal("failed to run db: %w", err)
	}

	defer func() {
		if err := cont.Terminate(ctx); err != nil {
			log.Fatal("failed to stop db: %w", err)
		}
	}()

	host, err = cont.ContainerIP(ctx)
	if err != nil {
		log.Fatal("failed to fetch db ip: %w", err)
	}

	code := m.Run()
	defer os.Exit(code)
}
