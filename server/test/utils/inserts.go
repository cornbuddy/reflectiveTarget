package utils

import (
	"database/sql"

	"github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
	"github.com/cornbuddy/reflectiveTarget/server/domain/entities"
	vo "github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
	_ "github.com/jackc/pgx/v5/stdlib"
)

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
	err := db.QueryRow(q, user.Username, user.Password.Hash).Scan(&user.ID)
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
