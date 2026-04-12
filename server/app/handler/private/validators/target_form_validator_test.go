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
		form     contracts.TargetForm
		wantForm contracts.TargetForm
		wantRes  bool
	}

	testCases := []testCase{{
		contracts.TargetForm{},
		contracts.TargetForm{},
		false,
	}}

	v := validators.TargetFormValidator{}
	for _, tc := range testCases {
		got := v.Validate(&tc.form)
		assert.Equal(t, tc.wantRes, got)
		assert.EqualExportedValues(t, tc.wantForm, tc.form)
	}
}
