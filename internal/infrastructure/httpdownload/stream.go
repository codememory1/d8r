package httpdownload

import (
	"context"
	"fmt"
	"net/http"

	"github.com/codememory1/d8r/internal/application/download"
	"github.com/codememory1/d8r/internal/infrastructure/storage"
)

// downloadStream downloads a resource without requiring its size in advance.
func (d *HttpDownloader) downloadStream(
	ctx context.Context,
	options download.Options,
	lifecycle download.Lifecycle,
) error {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		options.URL.String(),
		nil,
	)

	if err != nil {
		return err
	}

	if lifecycle.OnActiveRequestsChanged != nil {
		lifecycle.OnActiveRequestsChanged(1)

		defer lifecycle.OnActiveRequestsChanged(-1)
	}

	resp, err := d.client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	// Size may be unknown for a stream.
	var expectedSize *int64

	if options.Size != nil {
		size := options.Size.Int64()
		expectedSize = &size
	}

	writer, err := d.storage.CreateWriter(
		ctx,
		d.resolveFilename(options),
		d.options.BufferSize,
		expectedSize,
	)

	if err != nil {
		return err
	}

	defer writer.Close()

	written, err := writer.Write(ctx, resp.Body, storage.Lifecycle{
		OnProgress: lifecycle.OnProgress,
	})

	if err != nil {
		return err
	}

	if expectedSize != nil && written != *expectedSize {
		return fmt.Errorf(
			"download size mismatch: expected %d bytes, wrote %d bytes",
			*expectedSize,
			written,
		)
	}

	return nil
}
