package repositories

import (
	"cmp"
	"context"
	"database/sql"
	"errors"
	"slices"
	"strings"

	"go.uber.org/zap"

	aggr "github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
	vo "github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
	"github.com/cornbuddy/reflectiveTarget/server/infra/log"
)

type TargetRepo struct {
	DB *sql.DB
}

// creates or updates the target
func (r TargetRepo) Save(ctx context.Context, target *aggr.Target) error {
	log := log.Logger(ctx).With(zap.Stringer("target", target))

	log.Debug("going to start transaction")
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Error("failed to rollback transaction", zap.Error(err))
		}
	}()

	log.Debug("transaction started, going to save target")
	if err := r.saveTarget(ctx, tx, target); err != nil {
		return err
	}

	log.Debug("target saved, going to save questions")
	for i, question := range target.Questions {
		log := log.With(zap.Stringer("question", &question))
		log.Debug("going to save question")
		if err := r.saveQuestion(ctx, tx, target, &question); err != nil {
			return err
		}

		log.Debug("question is saved")
		target.Questions[i] = question
	}

	log.Debug("questions are saved, going to commit transaction")
	if err := tx.Commit(); err != nil {
		return err
	}

	log.Debug("target is committed, sorting questions")
	slices.SortFunc(target.Questions, func(a, b vo.Question) int {
		return cmp.Compare(a.ID, b.ID)
	})

	log.Debug("questions are sorted")

	return nil
}

// returns list of hollow (without nested fields) targets
func (r TargetRepo) ListTargetsOfUser(
	ctx context.Context, ownerID vo.ID,
) (aggr.Targets, error) {
	q := strings.Join([]string{
		"SELECT t.name, t.id",
		"FROM targets AS t",
		"WHERE t.owner_id = $1",
	}, "\n")
	rows, err := r.DB.QueryContext(ctx, q, ownerID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var res aggr.Targets
	for rows.Next() {
		t := aggr.Target{}
		if err := rows.Scan(&t.Name, &t.ID); err != nil {
			return nil, err
		}

		res = append(res, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (r TargetRepo) Get(ctx context.Context, id vo.ID) (*aggr.Target, error) {
	target, err := r.getTarget(ctx, id)
	if err != nil {
		return nil, err
	} else if target == nil {
		return nil, nil
	}

	questions, err := r.getQuestions(ctx, id)
	if err != nil {
		return nil, err
	}

	shots, err := r.getShots(ctx, id)
	if err != nil {
		return nil, err
	}

	target.Questions = questions
	target.Shots = shots

	return target, nil
}

func (r TargetRepo) getTarget(
	ctx context.Context, id vo.ID,
) (*aggr.Target, error) {
	res := &aggr.Target{}
	q := strings.Join([]string{
		"SELECT t.id, t.name, u.username, u.id, u.hashed_password",
		"FROM targets AS t",
		"JOIN users AS u ON t.owner_id = u.id",
		"WHERE t.id = $1",
	}, "\n")

	if err := r.DB.QueryRowContext(ctx, q, id).Scan(
		&res.ID, &res.Name, &res.Owner.Username, &res.Owner.ID,
		&res.Owner.Password.Hash,
	); errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return res, nil
}

const insertTarget = `
INSERT INTO targets (name, owner_id)
VALUES ($1, $2::integer)
RETURNING id`

const updateTarget = `
INSERT INTO targets AS t (id, name, owner_id)
VALUES ($1::integer, $2, $3::integer)
ON CONFLICT (id) DO UPDATE SET
name = EXCLUDED.name, owner_id = EXCLUDED.owner_id
WHERE (t.name, t.owner_id)
IS DISTINCT FROM (EXCLUDED.name, EXCLUDED.owner_id)
RETURNING id`

func (r TargetRepo) saveTarget(
	ctx context.Context, tx *sql.Tx, target *aggr.Target,
) error {
	log := log.Logger(ctx).With(zap.Stringer("target", target))

	var row *sql.Row
	name := target.Name
	ownerID := target.Owner.ID
	if target.ID == 0 {
		log.Debug("going to insert target")
		row = tx.QueryRowContext(ctx, insertTarget, name, ownerID)
		log.Debug("insert target query is executed")
	} else {
		log.Debug("going to update target")
		row = tx.QueryRowContext(ctx, updateTarget, target.ID, name, ownerID)
		log.Debug("update target query is executed")
	}

	if err := row.Scan(&target.ID); errors.Is(err, sql.ErrNoRows) {
		// kinda expected, this means update request did not update
		// anything
		return nil
	} else {
		return err
	}
}

func (r TargetRepo) getQuestions(
	ctx context.Context, targetId vo.ID,
) (vo.Questions, error) {
	q := strings.Join([]string{
		"SELECT q.id, q.text",
		"FROM questions AS q",
		"WHERE q.target_id = $1",
		"ORDER BY q.id",
	}, "\n")
	rows, err := r.DB.QueryContext(ctx, q, targetId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	questions := vo.Questions{}
	for rows.Next() {
		q := vo.Question{}
		if err := rows.Scan(&q.ID, &q.Text); err != nil {
			return nil, err
		}

		questions = append(questions, q)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return questions, nil
}

func (r TargetRepo) getShots(
	ctx context.Context, targetId vo.ID,
) (vo.Shots, error) {
	q := strings.Join([]string{
		"SELECT s.x, s.y",
		"FROM shots AS s",
		"WHERE s.target_id = $1",
	}, "\n")
	rows, err := r.DB.QueryContext(ctx, q, targetId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	shots := vo.Shots{}
	for rows.Next() {
		s := vo.Shot{}
		if err := rows.Scan(&s.X, &s.Y); err != nil {
			return nil, err
		}

		shots = append(shots, s)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return shots, nil
}

const insertQuestion = `
INSERT INTO questions (text, target_id)
VALUES ($1, $2::integer)
RETURNING id`

const updateQuestion = `
INSERT INTO questions AS q (id, text, target_id)
VALUES ($1::integer, $2, $3::integer)
ON CONFLICT (id) DO UPDATE SET
text = EXCLUDED.text, target_id = EXCLUDED.target_id
WHERE (q.text, q.target_id)
IS DISTINCT FROM (EXCLUDED.text, EXCLUDED.target_id)
RETURNING id`

func (r TargetRepo) saveQuestion(
	ctx context.Context, tx *sql.Tx, target *aggr.Target, question *vo.Question,
) error {
	log := log.Logger(ctx).With(zap.Stringer("question", question))

	var row *sql.Row
	text := question.Text
	if question.ID == 0 {
		log.Debug("going to insert question")
		row = tx.QueryRowContext(ctx, insertQuestion, text, target.ID)
		log.Debug("insert question query is executed")
	} else {
		log.Debug("going to update question")
		row = tx.QueryRowContext(ctx, updateQuestion, question.ID, text, target.ID)
		log.Debug("update question query is executed")
	}

	if err := row.Scan(&question.ID); errors.Is(err, sql.ErrNoRows) {
		// kinda expected, this means update request did not update
		// anything
		return nil
	} else {
		return err
	}
}
