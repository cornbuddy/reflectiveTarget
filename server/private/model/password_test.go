package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordHashingShouldNotBeDetermenistic(t *testing.T) {
	plaintext := "kek"
	pwd1, err := NewPassword(plaintext)
	assert.NoError(t, err)

	pwd2, err := NewPassword(plaintext)
	assert.NoError(t, err)

	assert.NotEqual(t, pwd1.Hash, pwd2.Hash)
}
