package render_test

import (
	"errors"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
)

var (
	err1 = errors.New("kek1")
	err2 = errors.New("kek2")
	errs = contracts.Errors{err1, err2}
)
