package daos

import (
	"database/sql"

	"github.com/cornbuddy/reflectiveTarget/server/private/model"
)

type ShotsDao struct {
	*sql.DB
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
