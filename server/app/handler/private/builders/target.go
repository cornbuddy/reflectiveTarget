package builders

import (
	"context"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/validators"
	"github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
	"github.com/cornbuddy/reflectiveTarget/server/domain/entities"
	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
)

type TargetBuilder struct {
	Validator validators.TargetFormValidator
}

// builds target from the form. returns nil if form is invalid. validates form
// as a side effect
func (b TargetBuilder) Target(
	ctx context.Context, form *contracts.TargetForm, ownerID valueobjects.ID,
	targetID valueobjects.ID,
) (*aggregations.Target, error) {
	if valid, err := b.Validator.Validate(ctx, form, ownerID, targetID); !valid {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	questions := make(valueobjects.Questions, 0, len(form.Questions))
	for _, field := range form.Questions {
		q := valueobjects.Question{ID: field.ID, Text: field.Value}
		questions = append(questions, q)
	}

	return &aggregations.Target{
		ID:        targetID,
		Name:      form.Name.Value,
		Owner:     entities.User{ID: ownerID},
		Questions: questions,
	}, nil
}
