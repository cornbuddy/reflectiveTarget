package handlers

import (
	"github.com/cornbuddy/reflectiveTarget/server/app/forms"
	"github.com/cornbuddy/reflectiveTarget/server/infra/daos"
)

type LoginFormValidator struct {
	daos.UserDao
}

func (v LoginFormValidator) Validate(form *forms.LoginForm) bool {
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

	user, err := v.UserDao.Find(form.Username.Value)
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

	valid, err := user.Password.Verify(form.Password.Value)
	if err != nil {
		form.Password.AddError(err)
		return false
	} else if !valid {
		form.Password.AddError(ErrWrongPassword)
	}

	return form.Username.IsValid() && form.Password.IsValid()
}
