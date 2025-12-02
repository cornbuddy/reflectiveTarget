package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cornbuddy/reflectiveTarget/server/app/formdata"
)

func TestSignupFormValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		form formdata.SignupForm
		want ValidationResult
	}

	testCases := []testCase{{}}

	validator := SignupFormValidator{}
	for _, tc := range testCases {
		got := validator.Validate(tc.form)
		assert.Equal(t, tc.want, got)
	}

}
