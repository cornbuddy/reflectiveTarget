package errors

import (
	"errors"
)

var ErrNotFound = errors.New("not found")
var ErrShotBadCoordinate = errors.New("bad coordinate for shot")
