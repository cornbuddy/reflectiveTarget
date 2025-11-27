package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggerInstanceShouldBeTheSameBetweenCalls(t *testing.T) {
	t.Parallel()

	l1 := Log
	assert.NotEmpty(t, l1)

	l2 := Log
	assert.NotEmpty(t, l2)
	assert.Same(t, l1, l2)
}
