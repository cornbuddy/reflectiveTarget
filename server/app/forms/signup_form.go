package forms

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
		Username:     Field{form.Get("username"), noErrs},
		Password:     Field{form.Get("password"), noErrs},
		Confirmation: Field{form.Get("confirmation"), noErrs},
	}
}

var noErrs Errors
