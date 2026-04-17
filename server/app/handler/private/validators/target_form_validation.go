package validators

import (
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
)

const (
	MaxNumOfQuestions = 10
	MaxTargetNameLen  = 128
	MaxQuestionLen    = 128
)

type TargetFormValidator struct{}

func (v *TargetFormValidator) Validate(form *contracts.TargetForm) bool {
	if len(form.Name.Value) == 0 {
		form.Name.AddError(ErrEmpty)
	}

	if len(form.Name.Value) > MaxTargetNameLen {
		form.Name.AddError(ErrTooLongTargetName)
	}

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
	}

	if len(form.Questions) == 0 {
		form.Questions = []contracts.Field{{
			Errors: contracts.Errors{ErrEmpty},
		}}
	}

	return form.Name.IsValid() && form.Questions.AreValid()
}
