package utils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"

	testdb "github.com/cornbuddy/reflectiveTarget/server/test/db"
)

func TestInitDatabase(t *testing.T) {
	ctx := context.TODO()
	db, err := testdb.Setup(ctx)
	assert.NoError(t, err)

	t.Cleanup(func() {
		err := testcontainers.TerminateContainer(db)
		assert.NoError(t, err)
	})
}
