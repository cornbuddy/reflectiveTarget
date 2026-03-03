package middlewares

import (
	"testing"
)

func TestRemoveTrailingSlash(t *testing.T) {
	t.Parallel()

	type testCase struct {
		desc string
		path string
	}

	handler := mw.RemoveTrailingSlash(emptyStub)
}
