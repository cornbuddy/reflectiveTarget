package validators

import (
	"fmt"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
)

const (
	MaxCoordinate = 100
	MinCoordinate = 0
)

type ShotsRequestValidator struct{}

func (v ShotsRequestValidator) Validate(
	shots contracts.ShotsRequest,
) ValidationResult {

	res := ValidationResult{}
	for _, shot := range shots.Shots {
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
