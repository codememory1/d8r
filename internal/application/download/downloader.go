package download

import (
	"context"
)

// Downloader defines a contract for downloading resources.
type Downloader interface {
	// Download downloads a resource using the specified options.
	Download(ctx context.Context, options Options) error
}
