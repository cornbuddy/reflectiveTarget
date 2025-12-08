package handlers

import (
	"errors"
	"fmt"
)

var ErrEmpty = errors.New("field cannot be empty")

var ErrUserAlreadyExists = errors.New("user already exists")

var ErrPasswordsShouldMatch = errors.New("passwords should match")
var ErrPasswordDoesntContainSpecialChar = errors.New(
	"password should contain at least one special character",
)
var ErrPasswordDoesntContainDigit = errors.New(
	"password should contain at least one digit",
)
var ErrPasswordTooShort = fmt.Errorf(
	"password should contain at least %d characters", MinPasswordLength,
)
