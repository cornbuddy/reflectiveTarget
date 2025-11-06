package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	t.Parallel()

	config, err := Init()
	require.NoError(t, err)
	require.NotNil(t, config)
}
