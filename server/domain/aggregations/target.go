package aggregations

import (
	"github.com/cornbuddy/reflectiveTarget/server/domain/entities"
	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
)

type Target struct {
	Name  string
	Owner entities.User
	valueobjects.ID
	valueobjects.Shots
	valueobjects.Questions
}

type Targets []Target
