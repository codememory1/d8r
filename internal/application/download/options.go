package download

import (
	"time"

	"github.com/codememory1/d8r/internal/domain/valueobject"
)

// Options contains the resource metadata and request settings required
// to download a resource.
type Options struct {
	URL          valueobject.URL
	Headers      valueobject.Headers
	ContentType  *valueobject.ContentType
	Filename     *valueobject.Filename
	Size         *valueobject.ByteSize
	Strategy     valueobject.DownloadStrategy
	ETag         *string
	LastModified *time.Time
}
