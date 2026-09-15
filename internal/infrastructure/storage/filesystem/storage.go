package filesystem

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/codememory1/d8r/internal/application/storage"
)

type Storage struct {
	directory string
}

// NewStorage creates a file storage with the specified configuration.
func NewStorage(directory string) *Storage {
	return &Storage{
		directory: directory,
	}
}

// CreateWriter prepares the directory and file for writing the resource,
// then returns a new sequential writer with the specified buffer size.
func (s *Storage) CreateWriter(_ context.Context, path string, bufferSize int64, size *int64) (storage.Writer, error) {
	// Preparing the directory. The directory is created if it does not exist.
	if err := s.prepareDirectory(); err != nil {
		return nil, err
	}

	// We prepare the file into which the resource will be loaded.
	file, err := s.prepareFile(fmt.Sprintf("%s/%s", s.directory, path), size)

	if err != nil {
		return nil, err
	}

	return NewWriter(file, bufferSize)
}

// prepareDirectory prepares the directory where the files will be saved.
func (s *Storage) prepareDirectory() error {
	info, err := os.Stat(s.directory)

	if err != nil {
		// If the directory does not exist, it creates the entire chain. For example: /a/b/c
		if errors.Is(err, os.ErrNotExist) {
			return os.MkdirAll(s.directory, 0755)
		}

		return err
	}

	if info.IsDir() {
		return nil
	}

	return fmt.Errorf("the path %q is a file, not a directory", s.directory)
}

// prepareFile prepares a file—that is, it creates the file
// and allocates disk space for it if a size is specified.
func (s *Storage) prepareFile(path string, size *int64) (*os.File, error) {
	flags := os.O_CREATE | os.O_WRONLY

	// We add the clear flag if the size was not passed,
	// since incorrect data could result if the file already
	// contains content.
	if size == nil {
		flags |= os.O_TRUNC
	}

	file, err := os.OpenFile(path, flags, 0644)

	if err != nil {
		return nil, err
	}

	if size != nil {
		// Changes the file size to allocate disk space.
		if err := file.Truncate(*size); err != nil {
			file.Close()

			return nil, err
		}
	}

	return file, nil
}
