package download

import (
	"time"

	"github.com/codememory1/d8r/internal/domain/valueobject"
)

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
