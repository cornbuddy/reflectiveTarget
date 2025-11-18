package validator

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationResult(t *testing.T) {
	t.Parallel()

	type testCase struct {
		msg    string
		result ValidationResult
		valid  bool
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
		assert.Equal(t, tc.valid, tc.result.IsValid(), tc.msg)
		assert.Equal(t, !tc.valid, tc.result.IsInvalid(), tc.msg)
	}
}
