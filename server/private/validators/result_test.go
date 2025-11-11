package validators

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationResultIsValid(t *testing.T) {
	t.Parallel()

	type testCase struct {
		msg    string
		result ValidationResult
		want   bool
	}

	testCases := []testCase{{
		"should be valid by default",
		ValidationResult{},
		true,
	}, {
		"should be valid when errors are empty",
		ValidationResult{[]error{}},
		true,
	}, {
		"should not be valid when errors are not empty",
		ValidationResult{[]error{errors.New("kek")}},
		false,
	}}

	for _, tc := range testCases {
		got := tc.result.IsValid()
		assert.Equal(t, tc.want, got, tc.msg)
	}
}
