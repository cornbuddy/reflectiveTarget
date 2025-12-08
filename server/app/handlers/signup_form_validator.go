package handlers

import (
	"regexp"
	"strings"

	"github.com/cornbuddy/reflectiveTarget/server/app/forms"
	"github.com/cornbuddy/reflectiveTarget/server/infra/daos"
)

type SignupFormValidator struct {
	daos.UserDao
}

const MinPasswordLength = 8

var SpecialChars = []rune{
	'!', '@', '#', '$', '%', '%', '^', '&', '*', '(', ')', '-', '_', '+',
	'=', '`', '~', '[', ']', '{', '}', '\\', '|', '/', '"', '\'', '?', ',',
	'.', '<', '>',
}

var hasDigits = regexp.MustCompile(`^.+\d.+$`)

func (v SignupFormValidator) Validate(form *forms.SignupForm) bool {
	user, err := v.UserDao.Find(form.Username.Value)
	if err != nil {
		form.Username.AddError(err)
	} else if user != nil {
		form.Username.AddError(ErrUserAlreadyExists)
	}

	if len(form.Username.Value) == 0 {
		form.Username.AddError(ErrEmpty)
	}

	if len(form.Password.Value) < MinPasswordLength {
		form.Password.AddError(ErrPasswordTooShort)
	}

	if !hasDigits.MatchString(form.Password.Value) {
		form.Password.AddError(ErrPasswordDoesntContainDigits)
	}

	if !strings.ContainsAny(form.Password.Value, string(SpecialChars)) {
		form.Password.AddError(ErrPasswordDoesntContainSpecialChars)
	}

	if form.Password.Value != form.Confirmation.Value {
		form.Confirmation.AddError(ErrPasswordsShouldMatch)
	}

	return form.Username.IsValid() &&
		form.Password.IsValid() &&
		form.Confirmation.IsValid()
}
