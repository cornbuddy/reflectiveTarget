package repositories

import (
	"context"
	"database/sql"
	"errors"

	aggr "github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
)

type TargetRepo struct {
	*sql.DB
}

func (r TargetRepo) Save(ctx context.Context, target *aggr.Target) error {
	return errors.New("not implemented")
}

func (r TargetRepo) ListTargetNamesOfUser(
	ctx context.Context, username string,
) ([]string, error) {

	return nil, errors.New("not implemented")
}
