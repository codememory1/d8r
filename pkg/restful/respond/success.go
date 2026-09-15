package respond

// SuccessBody represents a successful API response.
type SuccessBody[T any] struct {
	Data T     `json:"data"`
	Meta *Meta `json:"meta,omitempty"`
}

// NewSuccessBody creates a successful response body with the provided data.
func NewSuccessBody[T any](data T) SuccessBody[T] {
	return SuccessBody[T]{Data: data}
}

// WithCursorPagination attaches cursor pagination metadata to the response.
func (b SuccessBody[T]) WithCursorPagination(limit int, nextCursor *string) SuccessBody[T] {
	if b.Meta == nil {
		b.Meta = &Meta{}
	}

	b.Meta.Pagination = new(NewCursorPagination(limit, nextCursor))

	return b
}

// ResponseBody returns the serializable success response body.
func (b SuccessBody[T]) ResponseBody() any {
	return b
}
