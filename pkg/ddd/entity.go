package ddd

// Entity represents a domain object identified by a unique ID.
type Entity[ID any] interface {
	ID() ID
}
