package daos_test

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
	"github.com/cornbuddy/reflectiveTarget/server/infra/daos"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

var (
	ctx = context.TODO()

	db       *sql.DB
	cache    *redis.Client
	health   daos.HealthDao
	shotsDao daos.ShotsDao
	userDao  daos.UserDao

	target   aggregations.Target
	user     entities.User
	password valueobjects.Password
	shots    valueobjects.Shots
)

func TestMain(m *testing.M) {
	cleanUpDb, testDb, err := utils.SetupDB(ctx)
	if err != nil {
		log.Fatalf("failed to setup db: %v", err)
	}

	cleanUpCache, testCache, err := utils.SetupCache(ctx)
	if err != nil {
		log.Fatalf("failed to setup db: %v", err)
	}

	pwd, err := valueobjects.NewPassword("kek")
	if err != nil {
		log.Fatalf("failed to create password: %v", err)
	}

	cache = testCache
	db = testDb
	userDao = daos.UserDao{DB: db}
	shotsDao = daos.ShotsDao{DB: db}
	health = daos.HealthDao{db, cache}

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

	if err := fillDatabase(ctx, db, &user, &target, shots); err != nil {
		log.Fatalf("failed to fill db: %v", err)
	}

	code, err := utils.RunAndCleanup(ctx, m, cleanUpCache, cleanUpDb)
	if err != nil {
		log.Fatalf("failed to cleanup: %v", err)
	}

	os.Exit(code)
}

func fillDatabase(
	ctx context.Context, db *sql.DB, user *entities.User,
	target *aggregations.Target, shots valueobjects.Shots,
) error {
	if err := utils.InsertUser(ctx, db, user); err != nil {
		return err
	}

	target.Owner = *user
	if err := utils.InsertTarget(ctx, db, target); err != nil {
		return err
	}

	q := "INSERT INTO shots (x, y, target_id, shooter) " +
		"VALUES ($1, $2, $3, $4)"
	shooter := "i'm-a-shooter"
	for _, shot := range shots {
		rows, err := db.QueryContext(ctx, q, shot.X, shot.Y, target.ID, shooter)
		switch {
		case err != nil:
			return err
		case rows.Err() != nil:
			return rows.Err()
		}

		defer rows.Close()
	}

	return nil
}
