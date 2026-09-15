package storage

import (
	"context"
)

// Storage provides access to resource storage.
//
//go:generate mockgen -source=storage.go -destination=../../mocks/resource_storage.go -package=mocks
type Storage interface {
	CreateWriter(ctx context.Context, path string, bufferSize int64, size *int64) (Writer, error)
}
