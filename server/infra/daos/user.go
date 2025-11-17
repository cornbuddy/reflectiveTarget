package daos

import (
	"database/sql"

	"github.com/cornbuddy/reflectiveTarget/server/model/entities"
)

type UserDao struct {
	*sql.DB
}

func (dao UserDao) Save(user *entities.User) error {
	q := "INSERT INTO users(username, hashed_password) " +
		"VALUES($1, $2) " +
		"RETURNING id"
	err := dao.QueryRow(q, user.Username, user.Password.Hash).Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (dao UserDao) Find(username string) (*entities.User, error) {
	var user entities.User
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
