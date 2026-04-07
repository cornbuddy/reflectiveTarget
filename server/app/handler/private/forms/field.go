package forms

type Errors []error

type Field struct {
	Value string
	Errors
}

func (f *Field) AddError(err error) {
	f.Errors = append(f.Errors, err)
}

func (f *Field) IsValid() bool {
	return len(f.Errors) == 0
}

func (f *Field) IsInvalid() bool {
	return len(f.Errors) > 0
}
