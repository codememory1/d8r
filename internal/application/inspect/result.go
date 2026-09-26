package inspect

import (
	"time"

	"github.com/codememory1/d8r/internal/domain/valueobject"
)

// Result contains the metadata and download strategy discovered for a remote
// resource.
type Result struct {
	// EffectiveURL is the final resource URL after redirects.
	EffectiveURL valueobject.URL

	// ContentType is the resource media type, when provided by the server.
	ContentType *valueobject.ContentType

	// Filename is the resolved resource filename, when available.
	Filename *valueobject.Filename

	// Size is the total resource size, when it can be determined.
	Size *valueobject.ByteSize

	SupportsParallelDownload bool

	// ETag is the entity tag returned by the server, when available.
	ETag *string

	// LastModified is the server-provided modification time, when available.
	LastModified *time.Time
}
