package handlers

import (
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
			res.Errors = append(res.Errors, ErrShotBadCoordinate)
		}

		if shot.Y > MaxCoordinate || shot.Y < MinCoordinate {
			res.Errors = append(res.Errors, ErrShotBadCoordinate)
		}
	}

	return res
}
