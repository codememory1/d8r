package pagination

// Cursor identifies the last item of a page for cursor-based pagination.
type Cursor struct {
	LastID    string
	Timestamp int64
}
