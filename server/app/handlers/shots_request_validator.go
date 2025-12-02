package handlers

import (
	"fmt"

	"github.com/cornbuddy/reflectiveTarget/server/domain/errors"
)

const (
	MaxCoordinate = 100
	MinCoordinate = 0
)

type ShotsRequestValidator struct{}

func (v ShotsRequestValidator) Validate(shots ShotsRequest) ValidationResult {

	res := ValidationResult{}
	for _, shot := range shots.Shots {
		if shot.X > MaxCoordinate || shot.X < MinCoordinate {
			err := fmt.Errorf("%w: x", errors.ErrShotBadCoordinate)
			res.Errors = append(res.Errors, err)
		}

		if shot.Y > MaxCoordinate || shot.Y < MinCoordinate {
			err := fmt.Errorf("%w: y", errors.ErrShotBadCoordinate)
			res.Errors = append(res.Errors, err)
		}
	}

	return res
}
