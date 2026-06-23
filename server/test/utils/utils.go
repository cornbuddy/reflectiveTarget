package utils

import (
	"context"
	"encoding/json"
	"math/rand/v2"
	"testing"

	"golang.org/x/sync/errgroup"
)

// runs test suite and executes clean up in parallel. returns zero if tests are
// fine, and non-zeor othervise. returns error if any of cleanups was failed
func RunAndCleanup(
	ctx context.Context, m *testing.M, cleanups ...Cleanup,
) (int, error) {
	code := m.Run()
	if len(cleanups) == 0 {
		return code, nil
	}

	g, _ := errgroup.WithContext(ctx)
	for _, cleanup := range cleanups {
		g.Go(cleanup)
	}

	if err := g.Wait(); err != nil {
		return code, err
	}

	return code, nil
}

// returns full copy of the val
func DeepCopy[T any](val T) (*T, error) {
	bytes, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}

	var res T
	json.Unmarshal(bytes, &res)

	return &res, nil
}

// returns random string with given length
func MakeRandomString(length int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	result := make([]rune, length)
	for i := range result {
		result[i] = letters[rand.IntN(len(letters))]
	}

	return string(result)
}
