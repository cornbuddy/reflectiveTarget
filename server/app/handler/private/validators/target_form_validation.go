package validators

import (
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
)

const MaxNumOfQuestions = 10

type TargetFormValidator struct{}

func (v *TargetFormValidator) Validate(form *contracts.TargetForm) bool {
	return true
}
