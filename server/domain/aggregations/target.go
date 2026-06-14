package aggregations

import (
	"fmt"
	"strings"

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

func (t *Target) String() string {
	return fmt.Sprintf(strings.Join([]string{
		"ID: %d, Name: '%s',",
		"Owner: {%s},",
		"Questions: [],",
		"Shots (count): %d"}, " "),
		t.ID, t.Name, t.Owner.String(), len(t.Shots),
	)
}
