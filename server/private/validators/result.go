package validators

type ValidationResult struct {
	Errors []error
}

func (r ValidationResult) IsValid() bool {
	return len(r.Errors) == 0
}
