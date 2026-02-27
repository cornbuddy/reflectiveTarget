package repositories

import (
	"context"
	"database/sql"
	"errors"

	aggr "github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
	vo "github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
)

type TargetRepo struct {
	*sql.DB
}

func (r TargetRepo) Get(ctx context.Context, id vo.ID) (*aggr.Target, error) {
	return nil, nil
}

func (r TargetRepo) Save(ctx context.Context, target *aggr.Target) error {
	return errors.New("not implemented")
}

func (r TargetRepo) ListTargetNamesOfUser(
	ctx context.Context, username string,
) ([]string, error) {

	return nil, errors.New("not implemented")
}
