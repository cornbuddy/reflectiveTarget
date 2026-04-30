package builders

import (
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/validators"
	"github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
	"github.com/cornbuddy/reflectiveTarget/server/domain/entities"
	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
)

type TargetBuilder struct {
	Validator validators.TargetFormValidator
}

// builds target from the form. validates form as a side effect
func (b TargetBuilder) Target(
	form *contracts.TargetForm,
	ownerID valueobjects.ID,
) *aggregations.Target {

	if valid := b.Validator.Validate(form); !valid {
		return nil
	}

	questions := make(valueobjects.Questions, 0, len(form.Questions))
	for _, field := range form.Questions {
		q := valueobjects.Question{Text: field.Value}
		questions = append(questions, q)
	}

	return &aggregations.Target{
		Name:      form.Name.Value,
		Owner:     entities.User{ID: ownerID},
		Questions: questions,
	}
}
