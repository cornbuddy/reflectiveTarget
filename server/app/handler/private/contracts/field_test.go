package contracts_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
)

func TestFieldsMethods(t *testing.T) {
	t.Parallel()

	fields := contracts.Fields{{}}
	assert.True(t, fields.AreValid())
	assert.False(t, fields.AreInvalid())

	fields = contracts.Fields{{Errors: contracts.Errors{errors.New("kek")}}}
	assert.False(t, fields.AreValid())
	assert.True(t, fields.AreInvalid())
}

func TestFieldMethods(t *testing.T) {
	t.Parallel()

	field := contracts.Field{}
	assert.Empty(t, field.Errors)
	assert.True(t, field.IsValid())
	assert.False(t, field.IsInvalid())

	err := errors.New("kek")
	field.AddError(err)
	assert.Contains(t, field.Errors, err)
	assert.Len(t, field.Errors, 1)
	assert.False(t, field.IsValid())
	assert.True(t, field.IsInvalid())
}
