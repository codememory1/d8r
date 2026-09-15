package inspect

import (
	"time"

	"github.com/codememory1/d8r/internal/domain/valueobject"
)

type Result struct {
	EffectiveURL     valueobject.URL
	ContentType      *valueobject.ContentType
	Filename         *valueobject.Filename
	Size             *valueobject.ByteSize
	DownloadStrategy valueobject.DownloadStrategy
	ETag             *string
	LastModified     *time.Time
}
