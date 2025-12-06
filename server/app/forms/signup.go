package forms

import (
	"net/url"
)

type SignupForm struct {
	Username        Field
	Password        Field
	ConfirmPassword Field
}

func NewSignupForm(form url.Values) SignupForm {
	return SignupForm{
		Username:        Field{form.Get("username"), noErrs},
		Password:        Field{form.Get("password"), noErrs},
		ConfirmPassword: Field{form.Get("repeated-password"), noErrs},
	}
}

var noErrs Errors
