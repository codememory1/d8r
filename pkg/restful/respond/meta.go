package respond

// Meta contains additional metadata associated with a successful response.
type Meta struct {
	Pagination *CursorPagination `json:"pagination,omitempty"`
}

// CursorPagination contains metadata for cursor-based pagination.
type CursorPagination struct {
	Limit      int     `json:"limit"`
	NextCursor *string `json:"next_cursor,omitempty"`
}

// NewMeta creates response metadata with pagination information.
func NewMeta(pagination *CursorPagination) Meta {
	return Meta{pagination}
}

// NewCursorPagination creates cursor-based pagination metadata.
func NewCursorPagination(limit int, nextCursor *string) CursorPagination {
	return CursorPagination{limit, nextCursor}
}
