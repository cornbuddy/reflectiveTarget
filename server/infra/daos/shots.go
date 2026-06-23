package daos

import (
	"context"
	"database/sql"

	vo "github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
)

type ShotsDao struct {
	DB *sql.DB
}

func (d ShotsDao) List(ctx context.Context, targetID vo.ID) (vo.Shots, error) {
	q := "SELECT x, y FROM shots WHERE target_id = $1"
	rows, err := d.DB.QueryContext(ctx, q, targetID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	shots := vo.Shots{}
	for rows.Next() {
		shot := vo.Shot{}
		if err := rows.Scan(&shot.X, &shot.Y); err != nil {
			return nil, err
		}

		shots = append(shots, shot)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return shots, nil
}

func (d ShotsDao) Save(
	ctx context.Context, shooter string, targetID vo.ID, shots vo.Shots,
) error {
	var err error
	var rows *sql.Rows
	// TODO: use copy protocol https://github.com/jackc/pgx/discussions/1545
	for _, shot := range shots {
		q := "INSERT INTO shots (x, y, target_id, shooter) " +
			"VALUES ($1, $2, $3, $4)"
		x := shot.X
		y := shot.Y
		rows, err = d.DB.QueryContext(ctx, q, x, y, targetID, shooter)
		if err != nil {
			return err
		}

		defer rows.Close()

		if err := rows.Err(); err != nil {
			return err
		}
	}

	return nil
}
