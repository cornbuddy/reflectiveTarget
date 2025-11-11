package validators

import (
	"github.com/cornbuddy/reflectiveTarget/server/private/model"
)

const (
	MaxCoordinate = 100
	MinCoordinate = 0
)

type ShotsRequestValidator struct{}

func (v ShotsRequestValidator) Validate(
	req model.ShotsRequest,
) ValidationResult {

	res := ValidationResult{}
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
