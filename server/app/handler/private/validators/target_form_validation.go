package validators

import (
	"errors"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
)

const MaxNumOfQuestions = 10

type TargetFormValidator struct{}

func (v *TargetFormValidator) Validate(
	form *contracts.TargetForm,
) ValidationResult {

	return ValidationResult{Errors: []error{errors.New("not implemented")}}
}
