package repositories

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	aggr "github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
	vo "github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
)

type TargetRepo struct {
	*sql.DB
}

func (r TargetRepo) Save(ctx context.Context, target *aggr.Target) error {
	return errors.New("not implemented")
}

// returns list of hollow (without nested fields) targets
func (r TargetRepo) ListTargetsOfUser(
	ctx context.Context, ownerID valueobjects.ID,
) (aggr.Targets, error) {

	q := strings.Join([]string{
		"SELECT t.name, t.id",
		"FROM targets AS t",
		"WHERE t.owner_id = $1",
	}, "\n")
	rows, err := r.QueryContext(ctx, q, ownerID)
	if err != nil {
		return nil, err
	}

	var res aggr.Targets
	for rows.Next() {
		t := aggr.Target{}
		if err := rows.Scan(&t.Name, &t.ID); err != nil {
			return nil, err
		}

		res = append(res, t)
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

	if err := r.QueryRowContext(ctx, q, id).Scan(
		&res.ID, &res.Name, &res.Owner.Username, &res.Owner.ID,
		&res.Owner.Password.Hash,
	); errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return res, nil
}

func (r TargetRepo) getQuestions(
	ctx context.Context, targetId vo.ID,
) (vo.Questions, error) {

	q := strings.Join([]string{
		"SELECT q.id, q.text",
		"FROM questions AS q",
		"WHERE q.target_id = $1",
	}, "\n")
	rows, err := r.QueryContext(ctx, q, targetId)
	if err != nil {
		return nil, err
	}

	questions := vo.Questions{}
	for rows.Next() {
		q := vo.Question{}
		if err := rows.Scan(&q.ID, &q.Text); err != nil {
			return nil, err
		}

		questions = append(questions, q)
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
	rows, err := r.QueryContext(ctx, q, targetId)
	if err != nil {
		return nil, err
	}

	shots := vo.Shots{}
	for rows.Next() {
		s := vo.Shot{}
		if err := rows.Scan(&s.X, &s.Y); err != nil {
			return nil, err
		}

		shots = append(shots, s)
	}

	return shots, nil
}
