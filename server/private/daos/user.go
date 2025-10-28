package daos

import (
	"database/sql"

	"github.com/cornbuddy/reflectiveTarget/server/private/dsl"
)

type UserDao struct {
	*sql.DB
}

func (dao UserDao) Find(username string) (*dsl.User, error) {
	var user dsl.User
	query := "SELECT id, username FROM users WHERE username = $1"
	err := dao.DB.QueryRow(query, username).Scan(&user.ID, &user.Username)
	if err == sql.ErrNoRows {
		// kinda expected
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}
