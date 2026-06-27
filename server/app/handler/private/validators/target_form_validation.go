package validators

import (
	"context"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
	"github.com/cornbuddy/reflectiveTarget/server/infra/repositories"
)

const (
	MaxNumOfQuestions = 10
	MaxTargetNameLen  = 128
	MaxQuestionLen    = 128
)

type TargetFormValidator struct {
	repositories.TargetRepo
}

func (v *TargetFormValidator) Validate(
	ctx context.Context, form *contracts.TargetForm, ownerID valueobjects.ID,
	targetID valueobjects.ID,
) (bool, error) {
	targets, err := v.ListTargetsOfUser(ctx, ownerID)
	if err != nil {
		return false, err
	}

	for _, target := range targets {
		nameAlreadyTaken := form.Name.Value == target.Name && target.ID != targetID
		if nameAlreadyTaken {
			form.Name.AddError(ErrTargetAlreadyExists)

			break
		}
	}

	if len(form.Name.Value) == 0 {
		form.Name.AddError(ErrEmpty)
	}

	if len(form.Name.Value) > MaxTargetNameLen {
		form.Name.AddError(ErrTooLongTargetName)
	}

	texts := make(map[string]bool, len(form.Questions))
	for i, q := range form.Questions {
		if len(q.Value) == 0 {
			form.Questions[i].AddError(ErrEmpty)
		}

		if len(q.Value) > MaxQuestionLen {
			form.Questions[i].AddError(ErrTooLongQuestion)
		}

		if i >= MaxNumOfQuestions {
			form.Questions[i].AddError(ErrExcessiveQuestion)
		}

		if _, found := texts[q.Value]; found {
			form.Questions[i].AddError(ErrRepeatedQuestion)
		}

		texts[q.Value] = true
	}

	if len(form.Questions) == 0 {
		form.Questions = contracts.QuestionFields{{
			Field: contracts.Field{
				Errors: contracts.Errors{ErrEmpty},
			},
		}}
	}

	return form.Name.IsValid() && form.Questions.AreValid(), nil
}
