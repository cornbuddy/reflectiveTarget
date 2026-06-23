package valueobjects_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
)

func TestPasswordShouldVerify(t *testing.T) {
	t.Parallel()

	type testCase struct {
		password string
		verify   string
		want     bool
	}

	testCases := []testCase{{
		password: kek,
		verify:   "not-kek",
		want:     false,
	}, {
		password: kek,
		verify:   kek,
		want:     true,
	}}

	for _, tc := range testCases {
		pwd, err := valueobjects.NewPassword(tc.password)
		require.NoError(t, err)

		got, err := pwd.Verify(tc.verify)
		require.NoError(t, err)
		assert.Equal(t, tc.want, got)
	}
}

func TestPasswordHashingShouldNotBeDetermenistic(t *testing.T) {
	t.Parallel()

	plaintext := "kek"
	pwd1, err := valueobjects.NewPassword(plaintext)
	require.NoError(t, err)

	pwd2, err := valueobjects.NewPassword(plaintext)
	require.NoError(t, err)

	assert.NotEqual(t, pwd1.Hash, pwd2.Hash)
}
