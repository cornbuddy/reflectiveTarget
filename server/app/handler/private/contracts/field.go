package contracts

import (
	"errors"
	"fmt"
	"strings"

	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
)

type Errors []error

type Fields []Field

type Field struct {
	ID     valueobjects.ID
	Value  string
	Errors Errors
}

func (errs Errors) Error() string {
	if len(errs) == 0 {
		return ""
	}

	return errors.Join(errs...).Error()
}

func (fs Fields) String() string {
	return toString(fs)
}

func (fs Fields) AreValid() bool {
	return fieldsAreValid(fs)
}

func (fs Fields) AreInvalid() bool {
	return !fieldsAreValid(fs)
}

func (f Field) String() string {
	var errs strings.Builder
	for i, err := range f.Errors {
		fmt.Fprintf(&errs, "'%s'", err.Error())
		if i < len(f.Errors)-1 {
			fmt.Fprintf(&errs, ", ")
		}
	}

	return fmt.Sprintf(
		"{ID: %d, Value: '%s', Errors: [%s]}",
		f.ID, f.Value, errs.String(),
	)
}

func (f *Field) AddError(err error) {
	f.Errors = append(f.Errors, err)
}

func (f Field) IsValid() bool {
	return len(f.Errors) == 0
}

func (f Field) IsInvalid() bool {
	return !f.IsValid()
}
