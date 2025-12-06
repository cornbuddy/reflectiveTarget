package forms_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cornbuddy/reflectiveTarget/server/app/forms"
)

func TestFieldShouldAddError(t *testing.T) {
	t.Parallel()

	field := forms.Field{}
	assert.Empty(t, field.Errors)

	err := errors.New("kek")
	field.AddError(err)
	assert.Contains(t, field.Errors, err)
	assert.Len(t, field.Errors, 1)
}
