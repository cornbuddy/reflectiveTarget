package daos

import (
	"context"
	"database/sql"
	"log"
	"math/rand/v2"
	"os"
	"testing"

	"github.com/cornbuddy/reflectiveTarget/server/private/model"
	"github.com/cornbuddy/reflectiveTarget/server/private/utils"
	testdb "github.com/cornbuddy/reflectiveTarget/server/test/db"
)

var (
	ctx      context.Context
	db       *sql.DB
	userDao  UserDao
	shotsDao ShotsDao

	user     model.User
	password model.Password
	target   model.Target
	shots    model.Shots
)

func TestMain(m *testing.M) {
	ctx = context.TODO()

	t := &testing.T{}
	cleanup, testDb, err := testdb.SetupTestDb(ctx, t)
	if err != nil {
		log.Fatalf("failed to setup db: %v", err)
	}

	if err := utils.InitDatabase(testDb); err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	pwd, err := model.NewPassword("kek")
	if err != nil {
		log.Fatalf("failed to create password: %v", err)
	}

	db = testDb
	userDao = UserDao{DB: db}
	shotsDao = ShotsDao{DB: db}

	password = *pwd
	user = model.User{
		Password: password,
		Username: "daos-user",
	}
	target = model.Target{
		Name: "kek?",
	}
	shots = model.Shots{
		model.Shot{X: rand.IntN(101), Y: rand.IntN(101)},
		model.Shot{X: rand.IntN(101), Y: rand.IntN(101)},
	}

	if err := fillDatabase(db, &user, &target, shots); err != nil {
		log.Fatalf("failed to fill db: %v", err)
	}

	code := m.Run()
	defer os.Exit(code)

	if err := cleanup(); err != nil {
		log.Fatalf("failed to cleanup test suite: %v", err)
	}

}

func fillDatabase(
	db *sql.DB, user *model.User, target *model.Target, shots model.Shots,
) error {

	q := "INSERT INTO users(username, hashed_password) " +
		"VALUES($1, $2) " +
		"RETURNING id"
	err := db.QueryRow(q, user.Username, user.Password.Hash).Scan(&user.ID)
	if err != nil {
		return err
	}

	target.OwnerId = user.ID

	q = "INSERT INTO targets (name, owner_id) VALUES ($1, $2) RETURNING id"
	err = db.QueryRow(q, "kek?", user.ID).Scan(&target.ID)
	if err != nil {
		return err
	}

	q = "INSERT INTO shots (x, y, target_id, shooter) " +
		"VALUES ($1, $2, $3, $4)"
	shooter := "i'm-a-shooter"
	for _, shot := range shots {
		_, err = db.Query(q, shot.X, shot.Y, target.ID, shooter)
		if err != nil {
			return err
		}
	}

	return nil
}
