package handlers

import (
	"errors"
	"fmt"
	"strings"
)

var ErrEmpty = errors.New("field cannot be empty")

var ErrUserAlreadyExists = errors.New("user already exists")
var ErrUserDoesNotExists = errors.New("user does not exist")

var ErrWrongPassword = errors.New("wrong password")
var ErrPasswordsShouldMatch = errors.New("passwords should match")
var ErrPasswordDoesntContainSpecialChars = fmt.Errorf(
	"password should contain at least one special character: %s",
	runesToString(SpecialChars),
)
var ErrPasswordDoesntContainDigits = errors.New(
	"password should contain at least one digit",
)
var ErrPasswordTooShort = fmt.Errorf(
	"password should contain at least %d characters", MinPasswordLength,
)

func runesToString(runes []rune) string {
	var result []string
	for _, r := range runes {
		str := fmt.Sprintf("`%s`", string(r))
		result = append(result, str)
	}

	return strings.Join(result, ", ")
}
