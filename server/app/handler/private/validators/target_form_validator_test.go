package validators_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/validators"
)

func TestTargetFormValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		form contracts.TargetForm
		want validators.ValidationResult
	}

	testCases := []testCase{{
		contracts.TargetForm{},
		validators.ValidationResult{},
	}}

	v := validators.TargetFormValidator{}
	for _, tc := range testCases {
		got := v.Validate(&tc.form)
		assert.Equal(t, tc.want, got)
	}
}
