package inspect

import "context"

// Inspector retrieves metadata and determines the download capabilities of a
// remote resource.
type Inspector interface {
	// Inspect examines the resource using the provided URL and HTTP headers.
	Inspect(ctx context.Context, url string, headers map[string]string) (Result, error)
}
