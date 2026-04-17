package contracts

import (
	"net/url"
	"slices"
	"strings"
)

type TargetForm struct {
	Name      Field
	Questions Fields
}

func NewTargetForm(form url.Values) TargetForm {
	keys := slices.DeleteFunc(mapKeys(form), func(key string) bool {
		return !strings.Contains(key, "question_")
	})
	slices.Sort(keys)

	var questions []Field
	for _, key := range keys {
		q := form.Get(key)
		questions = append(questions, Field{Value: q})
	}

	return TargetForm{
		Name:      Field{Value: form.Get("name")},
		Questions: questions,
	}
}

func mapKeys(m url.Values) []string {
	res := make([]string, len(m))
	for k := range m {
		res = append(res, k)
	}

	return res
}
