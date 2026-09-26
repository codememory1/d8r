package domain

// ValidationError represents a domain value validation failure.
type ValidationError struct {
	reason string
	cause  error
}

// NewValidationError creates a validation error with the provided reason
// and underlying cause.
func NewValidationError(reason string, cause error) *ValidationError {
	return &ValidationError{reason, cause}
}

// Error returns the validation failure reason.
func (v *ValidationError) Error() string {
	return v.reason
}

// Unwrap returns the error that caused the validation failure.
func (v *ValidationError) Unwrap() error {
	return v.cause
}
