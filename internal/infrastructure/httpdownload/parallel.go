package httpdownload

import (
	"context"
	"fmt"
	"net/http"

	"github.com/codememory1/d8r/internal/application/download"
	"github.com/codememory1/d8r/internal/infrastructure/storage"
	"golang.org/x/sync/errgroup"
)

// downloadParallel Performs a parallel download of a resource by splitting it into ranges
// and downloading them simultaneously, subject to a limit on the number
// of parallel requests.
func (d *HttpDownloader) downloadParallel(ctx context.Context, options download.Options, lifecycle download.Lifecycle) error {
	filename := d.resolveFilename(options)

	// Creates a writer to which the downloaded data will be written.
	writer, err := d.storage.CreateWriter(
		ctx,
		filename,
		d.options.BufferSize,
		new(options.Size.Int64()),
	)

	if err != nil {
		return err
	}

	defer writer.Close()

	g, ctx := errgroup.WithContext(ctx)

	g.SetLimit(d.options.MaxParallelParts)

	for i := 0; i < d.options.RangeParts; i++ {
		// Calculates the HTTP byte range for a specific part.
		httpRange, err := ResolveHTTPRange(options.Size.Int64(), i, d.options.RangeParts)

		if err != nil {
			return err
		}

		g.Go(func() error {
			return d.downloadPart(ctx, options.URL.String(), writer, httpRange, lifecycle)
		})
	}

	return g.Wait()
}

// downloadPart downloads the specified part of the resource
// and writes the received data to the writer using
// buffered reading.
func (d *HttpDownloader) downloadPart(
	ctx context.Context,
	url string,
	writer storage.Writer,
	httpRange HTTPRange,
	lifecycle download.Lifecycle,
) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return err
	}

	req.Header.Set("Range", httpRange.HeaderValue())

	partSize := httpRange.End - httpRange.Start + 1

	return d.retry.Do(ctx, func() error {
		if lifecycle.OnActiveRequestsChanged != nil {
			lifecycle.OnActiveRequestsChanged(1)

			defer lifecycle.OnActiveRequestsChanged(-1)
		}

		resp, err := d.client.Do(req)

		if err != nil {
			return err
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusPartialContent {
			return fmt.Errorf("unexpected HTTP status: %s", resp.Status)
		}

		// Loads a portion of the resource's content into the writer.
		written, err := writer.WriteAt(ctx, resp.Body, httpRange.Start, storage.Lifecycle{
			OnProgress: lifecycle.OnProgress,
		})

		if err != nil {
			return err
		}

		// The number of written bytes is being checked.
		if written != partSize {
			return fmt.Errorf(
				"download part size mismatch: expected %d bytes, wrote %d bytes",
				partSize,
				written,
			)
		}

		return nil
	})
}
