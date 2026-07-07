package contracts

import (
	"fmt"
)

type QuestionField struct{ Field }

type QuestionFields []QuestionField

const (
	QuestionFieldPrefix       = "question"
	QuestionFieldValuePostfix = "value"
	QuestionFieldIDPostfix    = "id"
)

func (fs QuestionFields) AreValid() bool {
	return fieldsAreValid(fs)
}

func (fs QuestionFields) AreInvalid() bool {
	return !fieldsAreValid(fs)
}

func (fs QuestionFields) String() string {
	return toString(fs)
}

// returns input name of question text
func (f *QuestionField) NameValue(i int) string {
	return fmt.Sprintf(
		"%s_%d_%s",
		QuestionFieldPrefix, i, QuestionFieldValuePostfix,
	)
}

// returns input name of question ID
func (f *QuestionField) NameID(i int) string {
	return fmt.Sprintf(
		"%s_%d_%s",
		QuestionFieldPrefix, i, QuestionFieldIDPostfix,
	)
}
