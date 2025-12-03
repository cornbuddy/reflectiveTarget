package forms

type Errors []error

type Field struct {
	Value string
	Errors
}

func (f *Field) AddError(err error) {
	f.Errors = append(f.Errors, err)
}
