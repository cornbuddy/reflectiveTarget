package aggregations

import (
	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
)

type Target struct {
	Name string
	valueobjects.ID
	valueobjects.Shots
	valueobjects.Questions
}
