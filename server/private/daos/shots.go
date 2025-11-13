package daos

import (
	"database/sql"
	"errors"

	"github.com/cornbuddy/reflectiveTarget/server/private/model"
)

type ShotsDao struct {
	*sql.DB
}

func (d ShotsDao) List(targetID int) (model.Shots, error) {
	q := "SELECT * FROM targets WHERE id = $1"
	err := d.DB.QueryRow(q, targetID).Scan()
	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrNotFound
	}

	q = "SELECT x, y FROM shots WHERE target_id = $1"
	rows, err := d.DB.Query(q, targetID)
	if err != nil {
		return nil, err
	}

	shots := model.Shots{}
	for rows.Next() {
		shot := model.Shot{}
		if err := rows.Scan(&shot.X, &shot.Y); err != nil {
			return nil, err
		}

		shots = append(shots, shot)
	}

	return shots, nil
}

func (d ShotsDao) Save(shooter string, targetID int, shots model.Shots) error {
	q := "INSERT INTO shots (x, y, target_id, shooter) " +
		"VALUES ($1, $2, $3, $4)"
	for _, shot := range shots {
		_, err := d.DB.Query(q, shot.X, shot.Y, targetID, shooter)
		if err != nil {
			return err
		}

	}
	return nil
}
