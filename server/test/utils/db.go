package utils

import (
	"context"
	"database/sql"
	"time"

	"github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
	"github.com/cornbuddy/reflectiveTarget/server/domain/entities"
	vo "github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
	"github.com/cornbuddy/reflectiveTarget/server/infra/migrations"
	_ "github.com/jackc/pgx/v5/stdlib"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const DbImage = "postgres:18-alpine"

func InsertShots(
	db *sql.DB, shots vo.Shots, targetId vo.ID, shooter string,
) error {

	for i := range shots {
		err := InsertShot(db, &shots[i], targetId, shooter)
		if err != nil {
			return err
		}
	}

	return nil
}

func InsertShot(
	db *sql.DB, shot *vo.Shot, targetId vo.ID, shooter string,
) error {

	q := "INSERT INTO shots (x, y, target_id, shooter) " +
		"VALUES($1, $2, $3, $4) "
	_, err := db.Query(q, shot.X, shot.Y, targetId, shooter)
	if err != nil {
		return err
	}

	return nil
}

func InsertQuestions(
	db *sql.DB, questions vo.Questions, targetId vo.ID,
) error {

	for i := range questions {
		err := InsertQuestion(db, &questions[i], targetId)
		if err != nil {
			return err
		}
	}

	return nil
}

func InsertQuestion(db *sql.DB, question *vo.Question, targetId vo.ID) error {
	q := "INSERT INTO questions (text, target_id) VALUES($1, $2) " +
		"RETURNING id"
	err := db.QueryRow(q, question.Text, targetId).Scan(&question.ID)
	if err != nil {
		return err
	}

	return nil
}

func InsertUsers(db *sql.DB, users entities.Users) error {
	for i := range users {
		if err := InsertUser(db, &users[i]); err != nil {
			return err
		}
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

func InsertTargets(db *sql.DB, targets aggregations.Targets) error {
	for i := range targets {
		if err := InsertTarget(db, &targets[i]); err != nil {
			return err
		}
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

	id := target.ID
	if err := InsertQuestions(db, target.Questions, id); err != nil {
		return err
	}

	if err := InsertShots(db, target.Shots, id, "kek"); err != nil {
		return err
	}

	return nil
}

// starts and initializes database
func SetupDB(ctx context.Context) (Cleanup, *sql.DB, error) {
	cleanup, db, err := StartDB(ctx)
	if err != nil {
		return cleanup, nil, err
	}

	if _, err := db.Exec(migrations.InitQuery); err != nil {
		return cleanup, nil, err
	}

	return cleanup, db, nil
}

// starts testcontainer with database
func StartDB(ctx context.Context) (Cleanup, *sql.DB, error) {
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
