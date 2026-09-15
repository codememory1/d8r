package domain

type ValidationError struct {
	reason string
	cause  error
}

func NewValidationError(reason string, cause error) *ValidationError {
	return &ValidationError{reason, cause}
}

func (v *ValidationError) Error() string {
	return v.reason
}

func (v *ValidationError) Unwrap() error {
	return v.cause
}
