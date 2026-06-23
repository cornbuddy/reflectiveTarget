package daos

import (
	"context"
	"database/sql"
	"errors"

	"github.com/cornbuddy/reflectiveTarget/server/domain/entities"
)

type UserDao struct {
	DB *sql.DB
}

func (dao UserDao) Save(ctx context.Context, user *entities.User) error {
	query := "INSERT INTO users(username, hashed_password) " +
		"VALUES($1, $2) " +
		"RETURNING id"
	username := user.Username
	hash := user.Password.Hash
	err := dao.DB.QueryRowContext(ctx, query, username, hash).Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (dao UserDao) Find(
	ctx context.Context, username string,
) (*entities.User, error) {
	var user entities.User
	query := "SELECT id, username, hashed_password FROM users " +
		"WHERE username = $1"
	err := dao.DB.QueryRowContext(ctx, query, username).
		Scan(&user.ID, &user.Username, &user.Password.Hash)
	if errors.Is(err, sql.ErrNoRows) {
		// kinda expected
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}
