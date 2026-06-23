package middlewares_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

func TestContentTypeSetsContentType(t *testing.T) {
	t.Parallel()

	type testCase struct {
		desc  string
		key   string
		value string
	}

	testCases := []testCase{{
		"should handle well-known header",
		"Content-Type",
		"application/json",
	}, {
		"should handle custom header",
		"top",
		"kek",
	}}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			m := http.MethodPost
			h := mw.SetHeader(tc.key, tc.value)(emptyStub).ServeHTTP
			r, _, err := utils.MakeRequest("", m, "/", h, nil)
			require.NoError(t, err)
			assert.Equal(t, tc.value, r.Header.Get(tc.key))
		})
	}
}
