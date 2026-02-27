package utils

import (
	"context"
	"database/sql"
	"time"

	"github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
	"github.com/cornbuddy/reflectiveTarget/server/domain/entities"
	vo "github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
	_ "github.com/jackc/pgx/v5/stdlib"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const DbImage = "postgres:18-alpine"

func InsertQuestion(db *sql.DB, question *vo.Question, targetId vo.ID) error {
	q := "INSERT INTO questions (text, target_id) VALUES($1, $2) " +
		"RETURNING id"
	err := db.QueryRow(q, question.Text, targetId).Scan(&question.ID)
	if err != nil {
		return err
	}

	return nil
}

func InsertUser(db *sql.DB, user *entities.User) error {
	q := "INSERT INTO users (username, hashed_password) VALUES($1, $2) " +
		"RETURNING id"
	err := db.QueryRow(q, user.Username, user.Hash).Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil
}

func InsertTarget(db *sql.DB, target *aggregations.Target) error {
	q := "INSERT INTO targets (name, owner_id) VALUES($1, $2) " +
		"RETURNING id"
	err := db.QueryRow(q, target.Name, target.Owner.ID).Scan(&target.ID)
	if err != nil {
		return err
	}

	return nil
}

func SetupTestDb(ctx context.Context) (Cleanup, *sql.DB, error) {
	str := wait.ForLog("database system is ready to accept connections").
		WithOccurrence(2).
		WithStartupTimeout(5 * time.Second)
	opts := []tc.ContainerCustomizer{
		tc.WithWaitStrategy(str),
	}
	cont, err := postgres.Run(ctx, DbImage, opts...)

	cleanup := func() error {
		if err := tc.TerminateContainer(cont); err != nil {
			return err
		}

		return nil
	}

	if err != nil {
		return cleanup, nil, err
	}

	connStr, err := cont.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return cleanup, nil, err
	}

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return cleanup, nil, err
	}

	cleanup = func() error {
		if err := db.Close(); err != nil {
			return err
		}

		if err := tc.TerminateContainer(cont); err != nil {
			return err
		}

		return nil
	}

	return cleanup, db, nil
}
