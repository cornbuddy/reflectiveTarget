package validators

import (
	"context"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
	"github.com/cornbuddy/reflectiveTarget/server/infra/daos"
)

type LoginFormValidator struct {
	daos.UserDao
}

func (v LoginFormValidator) Validate(
	ctx context.Context, form *contracts.LoginForm,
) bool {
	emptyUsername := len(form.Username.Value) == 0
	if emptyUsername {
		form.Username.AddError(ErrEmpty)
	}

	emptyPassword := len(form.Password.Value) == 0
	if emptyPassword {
		form.Password.AddError(ErrEmpty)
	}

	if emptyUsername && emptyPassword {
		return false
	}

	user, err := v.Find(ctx, form.Username.Value)
	if err != nil {
		form.Username.AddError(err)

		return false
	} else if user == nil {
		form.Username.AddError(ErrUserDoesNotExists)
		// it doesn't make sense to validate password if there's no such
		// user. also, this check will fail anyway, because user
		// variable is nil
		return false
	}

	valid, err := user.Verify(form.Password.Value)
	if err != nil {
		form.Password.AddError(err)

		return false
	} else if !valid {
		form.Password.AddError(ErrWrongPassword)
	}

	return form.Username.IsValid() && form.Password.IsValid()
}
