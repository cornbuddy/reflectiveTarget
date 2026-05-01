package daos

import (
	"context"
	"database/sql"
	"log"
	"math/rand/v2"
	"os"
	"testing"

	"github.com/redis/go-redis/v9"

	"github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
	"github.com/cornbuddy/reflectiveTarget/server/domain/entities"
	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

var (
	ctx = context.TODO()

	db       *sql.DB
	cache    *redis.Client
	health   HealthDao
	store    SessionStore
	shotsDao ShotsDao
	userDao  UserDao

	target   aggregations.Target
	user     entities.User
	password valueobjects.Password
	shots    valueobjects.Shots
)

func TestMain(m *testing.M) {
	cleanUpDb, testDb, err := utils.StartDB(ctx)
	if err != nil {
		log.Fatalf("failed to setup db: %v", err)
	}

	defer func() {
		if err := cleanUpDb(); err != nil {
			log.Fatalf("failed to clean up db: %v", err)
		}
	}()

	cleanUpCache, testCache, err := utils.SetupCache(ctx)
	if err != nil {
		log.Fatalf("failed to setup db: %v", err)
	}

	defer func() {
		if err := cleanUpCache(); err != nil {
			log.Fatalf("failed to clean up cache: %v", err)
		}
	}()

	pwd, err := valueobjects.NewPassword("kek")
	if err != nil {
		log.Fatalf("failed to create password: %v", err)
	}

	cache = testCache
	db = testDb
	userDao = UserDao{DB: db}
	shotsDao = ShotsDao{DB: db}
	store = SessionStore{cache}
	health = HealthDao{db, cache}

	password = *pwd
	user = entities.User{
		Password: password,
		Username: "daos-user",
	}
	target = aggregations.Target{
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
	db *sql.DB, user *entities.User, target *aggregations.Target,
	shots valueobjects.Shots,
) error {

	if err := utils.InsertUser(db, user); err != nil {
		return err
	}

	target.Owner = *user
	if err := utils.InsertTarget(db, target); err != nil {
		return err
	}

	q := "INSERT INTO shots (x, y, target_id, shooter) " +
		"VALUES ($1, $2, $3, $4)"
	shooter := "i'm-a-shooter"
	for _, shot := range shots {
		_, err := db.Query(q, shot.X, shot.Y, target.ID, shooter)
		if err != nil {
			return err
		}
	}

	return nil
}
