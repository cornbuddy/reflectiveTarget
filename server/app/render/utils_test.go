package render_test

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var layoutMakrers = []string{"<!DOCTYPE html>", "</html>"}

func assertContainsTokens(t *testing.T, body io.Reader, tokens []string) {
	raw, err := io.ReadAll(body)
	require.NoError(t, err)

	strBody := string(raw)
	for _, token := range tokens {
		assert.Contains(t, strBody, token)
	}
}

func assertNotContainsTokens(t *testing.T, body io.Reader, tokens []string) {
	raw, err := io.ReadAll(body)
	require.NoError(t, err)

	strBody := string(raw)
	for _, token := range tokens {
		assert.NotContains(t, strBody, token)
	}
}
