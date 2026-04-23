package validators

import (
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
)

type ValidationResult struct {
	contracts.Errors
}

func (r ValidationResult) IsValid() bool {
	return len(r.Errors) == 0
}

func (r ValidationResult) IsInvalid() bool {
	return len(r.Errors) > 0
}
