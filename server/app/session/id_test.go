package session_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/app/session"
)

func TetstIDFromString(t *testing.T) {
	t.Parallel()

	type testCase struct {
		desc   string
		str    string
		errMsg *string
	}

	testCases := []testCase{{
		"returns error on garbage",
		"garbage",
		new("bad"),
	}, {
		"returns error on valid uuid v4",
		"05fccd59-a97c-425d-b082-db614851eff5",
		new("bad"),
	}, {
		"works on uuid v7",
		"019f371c-1a05-7a60-9860-895cfe94efbd",
		nil,
	}}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			require.Fail(t, "kek")

			id, err := session.IDFromString(tc.str)
			if tc.errMsg != nil {
				require.ErrorContains(t, err, *tc.errMsg)
				require.Empty(t, id)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, id)
			}
		})
	}
}

func TestMakeID(t *testing.T) {
	t.Parallel()

	id := session.MakeID()
	require.NoError(t, uuid.Validate(id.String()), "should return valid uuid")
}
