package daos

import (
	"context"
	"database/sql"
	"log"
	"math/rand/v2"
	"os"
	"testing"

	"github.com/cornbuddy/reflectiveTarget/server/domain/entities"
	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
	"github.com/cornbuddy/reflectiveTarget/server/infra/utils"
	testutils "github.com/cornbuddy/reflectiveTarget/server/test/utils"
	"github.com/redis/go-redis/v9"
)

var (
	ctx = context.TODO()

	db       *sql.DB
	cache    *redis.Client
	userDao  UserDao
	shotsDao ShotsDao
	store    SessionStore

	target   entities.Target
	user     entities.User
	password valueobjects.Password
	shots    valueobjects.Shots
)

func TestMain(m *testing.M) {
	cleanUpDb, testDb, err := testutils.SetupTestDb(ctx)
	if err != nil {
		log.Fatalf("failed to setup db: %v", err)
	}

	defer func() {
		if err := cleanUpDb(); err != nil {
			log.Fatalf("failed to clean up db: %v", err)
		}
	}()

	cleanUpCache, testCache, err := testutils.SetupCache(ctx)
	if err != nil {
		log.Fatalf("failed to setup db: %v", err)
	}

	defer func() {
		if err := cleanUpCache(); err != nil {
			log.Fatalf("failed to clean up cache: %v", err)
		}
	}()

	if err := utils.InitDatabase(testDb); err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	pwd, err := valueobjects.NewPassword("kek")
	if err != nil {
		log.Fatalf("failed to create password: %v", err)
	}

	cache = testCache
	db = testDb
	userDao = UserDao{DB: db}
	shotsDao = ShotsDao{DB: db}
	store = SessionStore{Ctx: ctx, Cache: cache}

	password = *pwd
	user = entities.User{
		Password: password,
		Username: "daos-user",
	}
	target = entities.Target{
		Name: "kek?",
	}
	shots = valueobjects.Shots{
		valueobjects.Shot{X: rand.IntN(101), Y: rand.IntN(101)},
		valueobjects.Shot{X: rand.IntN(101), Y: rand.IntN(101)},
	}

	if err := fillDatabase(db, &user, &target, shots); err != nil {
		log.Fatalf("failed to fill db: %v", err)
	}

	os.Exit(m.Run())
}

func fillDatabase(
	db *sql.DB, user *entities.User, target *entities.Target,
	shots valueobjects.Shots,
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
