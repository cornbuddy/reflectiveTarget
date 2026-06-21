package contracts_test

import (
	"errors"
)

const (
	formName     = "name"
	formPassword = "password"
	formUsername = "username"

	username = "kek"
	password = "pass"

	kek  = "kek"
	kek1 = "kek1"
	kek2 = "kek2"
	qek  = "kek?"
)

var (
	err  = errors.New(kek)
	err1 = errors.New(kek1)
	err2 = errors.New(kek2)
)
