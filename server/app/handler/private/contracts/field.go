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
	var res strings.Builder
	for i, f := range fs {
		fmt.Fprintf(&res, "%s", f.String())
		if i < len(fs)-1 {
			fmt.Fprintf(&res, ", ")
		}
	}

	return fmt.Sprintf("[%s]", res.String())
}

func (fs Fields) AreValid() bool {
	if len(fs) == 0 {
		return true
	}

	return all(fs, func(f Field) bool {
		return f.IsValid()
	})
}

func (fs Fields) AreInvalid() bool {
	return !fs.AreValid()
}

func (f *Field) String() string {
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

func (f *Field) IsValid() bool {
	return len(f.Errors) == 0
}

func (f *Field) IsInvalid() bool {
	return !f.IsValid()
}

func all[T any](ts []T, pred func(T) bool) bool {
	for _, t := range ts {
		if !pred(t) {
			return false
		}
	}

	return true
}
