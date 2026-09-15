package httpdownload

import (
	"context"
	"net/http"

	"github.com/codememory1/d8r/internal/application/download"
	"github.com/codememory1/d8r/internal/application/storage"
	"golang.org/x/sync/errgroup"
)

// downloadParallel Performs a parallel httpdownload of a resource by splitting it into ranges
// and downloading them simultaneously, subject to a limit on the number
// of parallel requests.
func (d *HttpDownloader) downloadParallel(ctx context.Context, options download.Options) error {
	filename := d.resolveFilename(options)

	// Creates a writer to which the downloaded data will be written.
	writer, err := d.storage.CreateWriter(
		ctx,
		filename,
		d.config.BufferSize.Bytes(),
		new(options.InspectionResult.Size.Int64()),
	)

	if err != nil {
		return err
	}

	defer writer.Close()

	g, ctx := errgroup.WithContext(ctx)

	g.SetLimit(d.config.MaxParallelParts)

	for i := 0; i < d.config.RangeParts; i++ {
		// Calculates the HTTP byte range for a specific part.
		httpRange, err := ResolveHTTPRange(options.InspectionResult.Size.Int64(), i, d.config.RangeParts)

		if err != nil {
			return err
		}

		g.Go(func() error {
			return d.downloadPart(ctx, options.InspectionResult.EffectiveURL.String(), writer, httpRange)
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
) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return err
	}

	req.Header.Set("Range", httpRange.HeaderValue())

	resp, err := d.client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// Loads a portion of the resource's content into the writer.
	_, err = writer.WriteAt(ctx, resp.Body, httpRange.Start)

	if err != nil {
		return err
	}

	return nil
}
