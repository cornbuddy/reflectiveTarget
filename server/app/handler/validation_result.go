package handler

type ValidationResult struct {
	Errors []error
}

func (r ValidationResult) IsValid() bool {
	return len(r.Errors) == 0
}

func (r ValidationResult) IsInvalid() bool {
	return len(r.Errors) > 0
}
