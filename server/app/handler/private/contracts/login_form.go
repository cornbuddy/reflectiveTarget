package contracts

import (
	"net/url"
)

type LoginForm struct {
	Username Field
	Password Field
}

func NewLoginForm(form url.Values) LoginForm {
	return LoginForm{
		Username: Field{Value: form.Get("username")},
		Password: Field{Value: form.Get("password")},
	}
}
