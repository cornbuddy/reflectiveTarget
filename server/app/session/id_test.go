package session_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/app/session"
)

func TestIDFromString(t *testing.T) {
	t.Parallel()

	type testCase struct {
		desc   string
		str    string
		errMsg *string
	}

	testCases := []testCase{{
		"returns error on garbage",
		"garbage",
		new("invalid UUID"),
	}, {
		"returns error on valid uuid v4",
		"69359037-9599-48e7-b8f2-48393c019135",
		new("bad version"),
	}, {
		"works on uuid v7",
		"019f371c-1a05-7a60-9860-895cfe94efbd",
		nil,
	}}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			id, err := session.IDFromString(tc.str)
			if tc.errMsg != nil {
				require.ErrorContains(t, err, *tc.errMsg)
				assert.Empty(t, id)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, id)
				assert.Equal(t, tc.str, id.String())
			}
		})
	}
}

func TestMakeID(t *testing.T) {
	t.Parallel()

	id := session.MakeID()
	require.NoError(t, uuid.Validate(id.String()), "should return valid uuid")
}
