package storage

import (
	"context"
	"io"
)

// Writer writes the resource content to the storage.
type Writer interface {
	// WriteAt writes data from reader, starting at the specified offset.
	WriteAt(ctx context.Context, reader io.Reader, offset int64) (int64, error)

	// Write sequentially writes data from the reader.
	Write(ctx context.Context, reader io.Reader) (int64, error)

	// Close releases the writer's resources.
	Close() error
}
