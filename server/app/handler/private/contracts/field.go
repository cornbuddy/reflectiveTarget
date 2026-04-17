package contracts

type Errors []error

type Fields []Field

type Field struct {
	Value string
	Errors
}

func (fs Fields) AreValid() bool {
	valid := false
	for _, f := range fs {
		if valid = f.IsValid(); valid {
			break
		}
	}

	return valid
}

func (fs Fields) AreInvalid() bool {
	invalid := false
	for _, f := range fs {
		if invalid = f.IsInvalid(); invalid {
			break
		}
	}

	return invalid
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
