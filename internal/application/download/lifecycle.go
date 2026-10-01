package download

// Lifecycle defines callbacks for download progress and request activity.
// Callbacks may be invoked concurrently during parallel downloads.
type Lifecycle struct {
	// OnProgress reports the number of bytes written by each write operation.
	OnProgress func(writtenBytes int64)

	// OnActiveRequestsChanged reports +1 when a request starts
	// and -1 when it finishes.
	OnActiveRequestsChanged func(delta int64)
}
