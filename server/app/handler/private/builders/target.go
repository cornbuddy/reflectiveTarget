package builders

import (
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/validators"
	"github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
)

type TargetBuilder struct {
	validator validators.TargetFormValidator
}

// builds target from the form. validates form as a side effect
func (b TargetBuilder) Target(form *contracts.TargetForm) *aggregations.Target {
	if valid := b.validator.Validate(form); !valid {
		return nil
	}

	return nil
}
