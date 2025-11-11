package daos

import (
	"database/sql"

	"github.com/cornbuddy/reflectiveTarget/server/private/model"
)

type UserDao struct {
	*sql.DB
}

func (dao UserDao) Save(user *model.User) error {
	q := "INSERT INTO users(username, hashed_password) " +
		"VALUES($1, $2) " +
		"RETURNING id"
	err := dao.QueryRow(q, user.Username, user.Password.Hash).Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (dao UserDao) Find(username string) (*model.User, error) {
	var user model.User
	query := "SELECT id, username, hashed_password FROM users " +
		"WHERE username = $1"
	err := dao.DB.QueryRow(query, username).
		Scan(&user.ID, &user.Username, &user.Password.Hash)
	if err == sql.ErrNoRows {
		// kinda expected
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}
