package repository

import "errors"

var (
	ErrNotFound               = errors.New("not found")
	ErrConcurrentModification = errors.New("entity was concurrently modified")
)
