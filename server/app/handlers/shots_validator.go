package handlers

import (
	"fmt"

	"github.com/cornbuddy/reflectiveTarget/server/infra/validator"
)

const (
	MaxCoordinate = 100
	MinCoordinate = 0
)

type ShotsRequestValidator struct{}

func (v ShotsRequestValidator) Validate(
	req ShotsRequest,
) validator.ValidationResult {

	res := validator.ValidationResult{}
	for _, shot := range req.Shots {
		if shot.X > MaxCoordinate || shot.X < MinCoordinate {
			err := fmt.Errorf("%w: x", ErrShotBadCoordinate)
			res.Errors = append(res.Errors, err)
		}

		if shot.Y > MaxCoordinate || shot.Y < MinCoordinate {
			err := fmt.Errorf("%w: y", ErrShotBadCoordinate)
			res.Errors = append(res.Errors, err)
		}
	}

	return res
}
