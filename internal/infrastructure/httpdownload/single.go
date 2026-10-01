package httpdownload

import (
	"context"
	"fmt"
	"net/http"

	"github.com/codememory1/d8r/internal/application/download"
	"github.com/codememory1/d8r/internal/infrastructure/storage"
)

// downloadSequential downloads the entire resource using a single HTTP request
// and writes it to storage sequentially.
func (d *HttpDownloader) downloadSequential(ctx context.Context, options download.Options, lifecycle download.Lifecycle) error {
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

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, options.URL.String(), nil)

	if err != nil {
		return err
	}

	return d.retry.Do(ctx, func() error {
		resp, err := d.client.Do(req)

		if err != nil {
			return err
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected HTTP status: %s", resp.Status)
		}

		written, err := writer.Write(ctx, resp.Body, storage.Lifecycle{
			OnProgress: lifecycle.OnProgress,
		})

		if err != nil {
			return err
		}

		// Verifies that the correct number of bytes has been written.
		if options.Size != nil && written != options.Size.Int64() {
			return fmt.Errorf(
				"download size mismatch: expected %d bytes, wrote %d bytes",
				options.Size.Int64(),
				written,
			)
		}

		return nil
	})
}
