package daos

import (
	"database/sql"

	"github.com/cornbuddy/reflectiveTarget/server/private/model"
)

type UserDao struct {
	*sql.DB
}

func (dao UserDao) Find(username string) (*model.User, error) {
	var user model.User
	query := "SELECT username FROM users WHERE username = $1"
	err := dao.DB.QueryRow(query, username).Scan(&user.Username)
	if err == sql.ErrNoRows {
		// kinda expected
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}
