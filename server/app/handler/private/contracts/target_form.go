package contracts

import (
	"cmp"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
)

var (
	ErrInvalidFieldValue = errors.New("bad value for field")

	questionValue = regexp.MustCompile(`^question_(\d+)_value$`)
)

type TargetForm struct {
	Name      Field
	Questions Fields
}

func NewTargetFormFromTarget(target aggregations.Target) TargetForm {
	questions := make(Fields, 0, len(target.Questions))
	for _, question := range target.Questions {
		field := Field{ID: question.ID, Value: question.Text}
		questions = append(questions, field)
	}

	return TargetForm{
		Name:      Field{Value: target.Name},
		Questions: questions,
	}
}

func NewTargetForm(form url.Values) (*TargetForm, error) {
	questions, err := parseQuestions(form)
	if err != nil {
		return nil, err
	}

	return &TargetForm{
		Name:      Field{Value: form.Get("name")},
		Questions: questions,
	}, nil
}

func (f *TargetForm) String() string {
	return fmt.Sprintf(
		"Name: %s, Questions: %s",
		f.Name.String(), f.Questions.String(),
	)
}

func parseQuestions(form url.Values) (Fields, error) {
	valueAttrs := slices.DeleteFunc(mapKeys(form), func(key string) bool {
		return !questionValue.MatchString(key)
	})
	slices.SortFunc(valueAttrs, func(a, b string) int {
		extractIndex := func(attr string) int {
			// ignoring errors since all attributes already
			// validated above
			ind, _ := strconv.Atoi(strings.Split(attr, "_")[1])

			return ind
		}

		ai := extractIndex(a)
		bi := extractIndex(b)

		return cmp.Compare(ai, bi)
	})

	var questions []Field
	for i, valueAttr := range valueAttrs {
		id := valueobjects.ID(0)
		idAttr := fmt.Sprintf("question_%d_id", i)
		rawID := form.Get(idAttr)
		if len(rawID) > 0 {
			intID, err := strconv.Atoi(rawID)
			if err != nil {
				return nil, fmt.Errorf(
					"%w: %s=%s (%w)",
					ErrInvalidFieldValue, idAttr, rawID, err,
				)
			}

			id = valueobjects.ID(intID)
		}

		questions = append(questions, Field{
			ID: id, Value: form.Get(valueAttr),
		})
	}

	return questions, nil
}

func mapKeys(m url.Values) []string {
	res := make([]string, len(m))
	for k := range m {
		res = append(res, k)
	}

	return res
}
