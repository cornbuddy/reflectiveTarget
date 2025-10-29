package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkNewPassword(b *testing.B) {
	b.Fail("not implemented")
}

func TestPasswordHashingShouldBeDetermenistic(t *testing.T) {
	plaintext := "kek"
	pwd1 := NewPassword(plaintext)
	pwd2 := NewPassword(plaintext)
	assert.Equal(t, pwd1.Hash, pwd2.Hash)
}
