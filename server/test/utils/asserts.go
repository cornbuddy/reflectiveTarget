package utils

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type assertFunc func(assert.TestingT, any, any, ...any) bool

func AssertContainsTokens(t *testing.T, body io.Reader, tokens []string) {
	t.Helper()

	assertTokens(t, body, tokens, assert.Contains)
}

func AssertNotContainsTokens(t *testing.T, body io.Reader, tokens []string) {
	t.Helper()

	assertTokens(t, body, tokens, assert.NotContains)
}

func assertTokens(
	t *testing.T, body io.Reader, tokens []string, asrt assertFunc,
) {
	t.Helper()

	raw, err := io.ReadAll(body)
	require.NoError(t, err)

	strBody := string(raw)
	for _, token := range tokens {
		asrt(t, strBody, token)
	}
}
