package contracts

type Errors []error

type Fields []Field

type Field struct {
	Value string
	Errors
}

func (fs Fields) AreValid() bool {
	if len(fs) == 0 {
		return true
	}

	return all(fs, func(f Field) bool {
		return f.IsValid()
	})
}

func (fs Fields) AreInvalid() bool {
	return !fs.AreValid()
}

func (f *Field) AddError(err error) {
	f.Errors = append(f.Errors, err)
}

func (f *Field) IsValid() bool {
	return len(f.Errors) == 0
}

func (f *Field) IsInvalid() bool {
	return !f.IsValid()
}

func all[T any](ts []T, pred func(T) bool) bool {
	for _, t := range ts {
		if !pred(t) {
			return false
		}
	}

	return true
}
