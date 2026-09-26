package filesystem

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
)

// WriterFunc writes a chunk of resource data and returns the number of bytes
// written.
type WriterFunc func(data []byte) (int, error)

// Writer writes resource data to an open filesystem file.
type Writer struct {
	file       *os.File
	bufferSize int64
}

// NewWriter creates a file Writer with a specified open file and buffer size.
func NewWriter(file *os.File, bufferSize int64) (*Writer, error) {
	if bufferSize <= 0 {
		return nil, fmt.Errorf("buffer size must be greater than zero")
	}

	return &Writer{file, bufferSize}, nil
}

// WriteAt reads data from reader and writes it to the file,
// starting at the specified offset.
func (w *Writer) WriteAt(_ context.Context, reader io.Reader, offset int64) (int64, error) {
	return w.copyTo(reader, func(data []byte) (int, error) {
		written, writeErr := w.file.WriteAt(data, offset)

		if writeErr != nil {
			return 0, writeErr
		}

		// Shift the write position by the number of bytes written,
		// so that the next chunk of data is written immediately after the previous one.
		offset += int64(written)

		return written, nil
	})
}

// Write sequentially reads data from the reader and writes it to the file.
func (w *Writer) Write(_ context.Context, reader io.Reader) (int64, error) {
	return w.copyTo(reader, func(data []byte) (int, error) {
		return w.file.Write(data)
	})
}

// Close closes the open file.
func (w *Writer) Close() error {
	return w.file.Close()
}

// copyTo reads data from reader and writes it using the provided function,
// calling onProgress after each successful write.
func (w *Writer) copyTo(reader io.Reader, writer WriterFunc) (int64, error) {
	buffer := make([]byte, w.bufferSize)

	var uploaded int64

	for {
		n, readErr := reader.Read(buffer)

		if n > 0 {
			written, writeErr := writer(buffer[:n])

			if writeErr != nil {
				return 0, writeErr
			}

			// Return an error if fewer bytes were written than read.
			if written != n {
				return 0, io.ErrShortWrite
			}

			// Increment the total number of successfully written bytes.
			uploaded += int64(written)
		}

		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				break
			}

			return 0, readErr
		}
	}

	return uploaded, nil
}
