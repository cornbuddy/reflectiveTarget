package validators

import (
	"context"
	"regexp"
	"strings"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
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

var hasDigits = regexp.MustCompile(`\d`)

func (v SignupFormValidator) Validate(
	ctx context.Context, form *contracts.SignupForm,
) (bool, error) {

	user, err := v.UserDao.Find(ctx, form.Username.Value)
	if err != nil {
		return false, err
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

	valid := form.Username.IsValid() &&
		form.Password.IsValid() &&
		form.Confirmation.IsValid()

	return valid, nil
}
