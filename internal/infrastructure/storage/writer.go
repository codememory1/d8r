package storage

import (
	"context"
	"io"
)

// Lifecycle contains optional callbacks for a storage write operation.
type Lifecycle struct {
	// OnProgress reports the cumulative bytes written by the current Write or WriteAt call.
	OnProgress func(writtenBytes int64)
}

// Writer writes the resource content to the storage.
type Writer interface {
	// WriteAt writes data from reader, starting at the specified offset.
	WriteAt(ctx context.Context, reader io.Reader, offset int64, lifecycle Lifecycle) (int64, error)

	// Write sequentially writes data from the reader.
	Write(ctx context.Context, reader io.Reader, lifecycle Lifecycle) (int64, error)

	// Close releases the writer's resources.
	Close() error
}
