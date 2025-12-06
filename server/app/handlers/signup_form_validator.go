package handlers

import (
	"github.com/cornbuddy/reflectiveTarget/server/app/forms"
	"github.com/cornbuddy/reflectiveTarget/server/infra/daos"
)

type SignupFormValidator struct {
	daos.UserDao
}

func (v SignupFormValidator) Validate(form *forms.SignupForm) {
	if len(form.Username.Value) == 0 {
		form.Username.AddError(ErrEmpty)
	}

	if len(form.Password.Value) == 0 {
		form.Password.AddError(ErrEmpty)
	}

	if len(form.Confirmation.Value) == 0 {
		form.Confirmation.AddError(ErrEmpty)
	}
}
