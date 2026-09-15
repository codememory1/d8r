package restful

// Error represents an HTTP error that can be safely returned to the client.
type Error struct {
	Status  int
	Message string
	Err     error
}

// NewError creates an HTTP error with a status, public message, and underlying cause.
func NewError(status int, message string, err error) Error {
	return Error{status, message, err}
}

// Error returns the public error message.
func (e Error) Error() string {
	return e.Message
}

// Unwrap returns the underlying error.
func (e Error) Unwrap() error {
	return e.Err
}
