package utils

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type assertFunc func(assert.TestingT, any, any, ...any) bool

// checks if response body contains substrings
func AssertContainsTokens(t *testing.T, r *http.Response, tokens []string) {
	t.Helper()

	assertTokens(t, r, tokens, assert.Contains)
}

// checks if response body does not contain substrings
func AssertNotContainsTokens(t *testing.T, r *http.Response, tokens []string) {
	t.Helper()

	assertTokens(t, r, tokens, assert.NotContains)
}

func assertTokens(
	t *testing.T, r *http.Response, tokens []string, asrt assertFunc,
) {
	t.Helper()

	raw, err := io.ReadAll(r.Body)
	require.NoError(t, err)
	require.NoError(t, r.Body.Close())

	r.Body = io.NopCloser(bytes.NewBuffer(raw))

	body := string(raw)
	for _, token := range tokens {
		asrt(t, body, token)
	}
}
