package daos

import (
	"context"
	"database/sql"
	"errors"

	myerrors "github.com/cornbuddy/reflectiveTarget/server/domain/errors"
	vo "github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
)

type ShotsDao struct {
	*sql.DB
}

func (d ShotsDao) List(ctx context.Context, targetID vo.ID) (vo.Shots, error) {
	q := "SELECT * FROM targets WHERE id = $1"
	err := d.DB.QueryRowContext(ctx, q, targetID).Scan()
	if errors.Is(err, sql.ErrNoRows) {
		return nil, myerrors.ErrNotFound
	}

	q = "SELECT x, y FROM shots WHERE target_id = $1"
	rows, err := d.DB.QueryContext(ctx, q, targetID)
	if err != nil {
		return nil, err
	}

	shots := vo.Shots{}
	for rows.Next() {
		shot := vo.Shot{}
		if err := rows.Scan(&shot.X, &shot.Y); err != nil {
			return nil, err
		}

		shots = append(shots, shot)
	}

	return shots, nil
}

func (d ShotsDao) Save(
	ctx context.Context, shooter string, targetID vo.ID, shots vo.Shots,
) error {

	q := "INSERT INTO shots (x, y, target_id, shooter) " +
		"VALUES ($1, $2, $3, $4)"
	for _, shot := range shots {
		x := shot.X
		y := shot.Y
		_, err := d.DB.QueryContext(ctx, q, x, y, targetID, shooter)
		if err != nil {
			return err
		}
	}

	return nil
}
