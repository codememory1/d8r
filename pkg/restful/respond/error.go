package respond

// ErrorBody represents an unsuccessful API response.
type ErrorBody struct {
	Message string `json:"message"`
}

// NewErrorBody creates a new error response body.
func NewErrorBody(message string) ErrorBody {
	return ErrorBody{message}
}

// ResponseBody returns the serializable error response body.
func (b ErrorBody) ResponseBody() any {
	return b
}
