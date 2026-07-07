package utils

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
	"github.com/cornbuddy/reflectiveTarget/server/domain/entities"
	vo "github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
)

func InsertShots(
	ctx context.Context, db *sql.DB, shots vo.Shots, targetId vo.ID,
	shooter string,
) error {
	for i := range shots {
		err := InsertShot(ctx, db, &shots[i], targetId, shooter)
		if err != nil {
			return err
		}
	}

	return nil
}

func InsertShot(
	ctx context.Context, db *sql.DB, shot *vo.Shot, targetId vo.ID,
	shooter string,
) error {
	q := "INSERT INTO shots (x, y, target_id, shooter) " +
		"VALUES($1, $2, $3, $4) "
	rows, err := db.QueryContext(ctx, q, shot.X, shot.Y, targetId, shooter)
	if err != nil {
		return err
	}

	defer rows.Close()

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

func InsertQuestions(
	ctx context.Context, db *sql.DB, questions vo.Questions, targetId vo.ID,
) error {
	for i := range questions {
		err := InsertQuestion(ctx, db, &questions[i], targetId)
		if err != nil {
			return err
		}
	}

	return nil
}

func InsertQuestion(
	ctx context.Context, db *sql.DB, question *vo.Question, targetId vo.ID,
) error {
	q := "INSERT INTO questions (text, target_id) VALUES($1, $2) " +
		"RETURNING id"
	err := db.QueryRowContext(ctx, q, question.Text, targetId).Scan(&question.ID)
	if err != nil {
		return err
	}

	return nil
}

func InsertUsers(ctx context.Context, db *sql.DB, users entities.Users) error {
	for i := range users {
		if err := InsertUser(ctx, db, &users[i]); err != nil {
			return err
		}
	}

	return nil
}

func InsertUser(ctx context.Context, db *sql.DB, user *entities.User) error {
	q := "INSERT INTO users (username, hashed_password) VALUES($1, $2) " +
		"RETURNING id"
	err := db.QueryRowContext(ctx, q, user.Username, user.Password.Hash).
		Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil
}

func InsertTargets(
	ctx context.Context, db *sql.DB, targets aggregations.Targets,
) error {
	for i := range targets {
		if err := InsertTarget(ctx, db, &targets[i]); err != nil {
			return err
		}
	}

	return nil
}

func InsertTarget(
	ctx context.Context, db *sql.DB, target *aggregations.Target,
) error {
	q := "INSERT INTO targets (name, owner_id) VALUES($1, $2) " +
		"RETURNING id"
	err := db.QueryRowContext(ctx, q, target.Name, target.Owner.ID).
		Scan(&target.ID)
	if err != nil {
		return err
	}

	id := target.ID
	if err := InsertQuestions(ctx, db, target.Questions, id); err != nil {
		return err
	}

	if err := InsertShots(ctx, db, target.Shots, id, "kek"); err != nil {
		return err
	}

	return nil
}
