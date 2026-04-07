package contracts

import (
	"net/url"
)

type SignupForm struct {
	Username     Field
	Password     Field
	Confirmation Field
}

func NewSignupForm(form url.Values) SignupForm {
	return SignupForm{
		Username:     Field{Value: form.Get("username")},
		Password:     Field{Value: form.Get("password")},
		Confirmation: Field{Value: form.Get("confirmation")},
	}
}
