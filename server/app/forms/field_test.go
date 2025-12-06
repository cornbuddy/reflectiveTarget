package forms

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldShouldAddError(t *testing.T) {
	t.Parallel()

	field := Field{}
	assert.Empty(t, field.Errors)

	err := errors.New("kek")
	field.AddError(err)
	assert.Contains(t, field.Errors, err)
	assert.Len(t, field.Errors, 1)
}
