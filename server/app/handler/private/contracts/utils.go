package contracts

import (
	"fmt"
	"strings"
)

type field interface {
	IsValid() bool
}

func fieldsAreValid[F field](fs []F) bool {
	if len(fs) == 0 {
		return true
	}

	return all(fs, func(f F) bool {
		return f.IsValid()
	})
}

func all[T any](ts []T, pred func(T) bool) bool {
	for _, t := range ts {
		if !pred(t) {
			return false
		}
	}

	return true
}

func toString[S fmt.Stringer](fs []S) string {
	var res strings.Builder
	for i, f := range fs {
		fmt.Fprintf(&res, "%s", f.String())
		if i < len(fs)-1 {
			fmt.Fprintf(&res, ", ")
		}
	}

	return fmt.Sprintf("[%s]", res.String())
}
