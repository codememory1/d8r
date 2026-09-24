package httpdownload

import (
	"context"
	"net/http"

	"github.com/codememory1/d8r/internal/application/download"
)

func (d *HttpDownloader) downloadSequential(ctx context.Context, options download.Options) error {
	filename := d.resolveFilename(options)

	// Creates a writer to which the downloaded data will be written.
	writer, err := d.storage.CreateWriter(
		ctx,
		filename,
		d.config.BufferSize.Bytes(),
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

	resp, err := d.client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	_, err = writer.Write(ctx, resp.Body)

	if err != nil {
		return err
	}

	return nil
}
