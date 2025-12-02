package handlers

import (
	"github.com/cornbuddy/reflectiveTarget/server/app/formdata"
)

type SignupFormValidator struct{}

func (v SignupFormValidator) Validate(
	form formdata.SignupForm,
) ValidationResult {

	return ValidationResult{}
}
