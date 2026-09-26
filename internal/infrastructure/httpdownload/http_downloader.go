package httpdownload

import (
	"context"
	"fmt"
	"mime"
	"net/http"

	"github.com/codememory1/d8r/internal/application/download"
	"github.com/codememory1/d8r/internal/application/storage"
	"github.com/codememory1/d8r/internal/domain/valueobject"
	"github.com/google/uuid"
)

// HttpDownloader downloads remote resources over HTTP and stores their contents.
type HttpDownloader struct {
	client  *http.Client
	storage storage.Storage
	options Options
}

// NewHttpDownloader creates a new HTTP downloader.
func NewHttpDownloader(client *http.Client, storage storage.Storage, options Options) *HttpDownloader {
	return &HttpDownloader{
		client:  client,
		storage: storage,
		options: options,
	}
}

// Download downloads a resource using the strategy selected during inspection.
func (d *HttpDownloader) Download(ctx context.Context, options download.Options) error {
	if options.Strategy.Equal(valueobject.SingleDownloadStrategy()) {
		return d.downloadSequential(ctx, options)
	}

	if options.Strategy.Equal(valueobject.StreamDownloadStrategy()) {
		return nil
	}

	if options.Strategy.Equal(valueobject.ParallelDownloadStrategy()) {
		return d.downloadParallel(ctx, options)
	}

	return nil
}

// resolveFilename returns an explicitly provided filename, a filename detected
// during inspection, or a generated UUID-based filename.
func (d *HttpDownloader) resolveFilename(options download.Options) string {
	if options.Filename != nil {
		return options.Filename.String()
	}

	if options.Filename != nil {
		return options.Filename.String()
	}

	contentType := options.ContentType

	if contentType != nil {
		ext, err := mime.ExtensionsByType(contentType.String())

		if err == nil && len(ext) > 0 {
			return fmt.Sprintf("%s.%s", uuid.NewString(), ext[0])
		}
	}

	return uuid.NewString()
}
