package formdata

import (
	"net/url"
)

type SignupForm struct {
	Username         string
	Password         string
	RepeatedPassword string
}

func NewSignupForm(form url.Values) SignupForm {
	return SignupForm{
		Username:         form.Get("username"),
		Password:         form.Get("password"),
		RepeatedPassword: form.Get("repeated-password"),
	}
}
