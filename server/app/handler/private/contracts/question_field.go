package contracts

import "fmt"

type QuestionField struct{ Field }

type QuestionFields []QuestionField

const (
	QuestionFieldPrefix       = "question"
	QuestionFieldValuePostfix = "value"
	QuestionFieldIDPostfix    = "id"
)

// returns input name of question text
func (f *QuestionField) NameValue() string {
	return fmt.Sprintf(
		"%s_%d_%s",
		QuestionFieldPrefix, f.ID, QuestionFieldValuePostfix,
	)
}

// returns input name of question ID
func (f *QuestionField) NameID() string {
	return fmt.Sprintf(
		"%s_%d_%s",
		QuestionFieldPrefix, f.ID, QuestionFieldIDPostfix,
	)
}
