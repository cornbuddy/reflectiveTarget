package contracts

import (
	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
)

type ShotsRequest struct {
	valueobjects.Shots `json:"shots"`
}

type ShotsResponse struct {
	valueobjects.Shots `json:"shots"`
}
