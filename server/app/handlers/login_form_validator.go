package handlers

import (
	"github.com/cornbuddy/reflectiveTarget/server/app/forms"
	"github.com/cornbuddy/reflectiveTarget/server/infra/daos"
)

type LoginFormValidator struct {
	daos.UserDao
}

func (v LoginFormValidator) Validate(form *forms.LoginForm) bool {
	return false
}
