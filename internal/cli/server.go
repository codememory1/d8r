package cli

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/codememory1/d8r/internal/bootstrap"
	"github.com/spf13/cobra"
)

// newServerCommand creates the command that starts the HTTP server.
func (rt *runtime) newServerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Run the server",
		RunE:  rt.runServer,
	}
}

// runServer starts the HTTP server and gracefully shuts it down when the
// command context is canceled.
func (rt *runtime) runServer(cmd *cobra.Command, _ []string) error {
	return rt.withApp(cmd, func(ctx context.Context, app *bootstrap.App) error {
		server, serveErr, err := serveHTTP(app)
		workerErr := runWorkers(ctx, app)

		if err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return shutdownHTTP(server)
		case err := <-serveErr:
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}

			return err
		case err := <-workerErr:
			if shutdownErr := shutdownHTTP(server); err != nil {
				return errors.Join(err, shutdownErr)
			}

			return err
		}
	})
}

// serveHTTP starts the HTTP server and returns a channel through which its
// terminal error is reported.
func serveHTTP(app *bootstrap.App) (*http.Server, <-chan error, error) {
	address := fmt.Sprintf("%s:%d", app.Config.HTTP.Address, app.Config.HTTP.Port)
	server := &http.Server{
		Addr:    address,
		Handler: app.Router(),
	}
	listener, err := net.Listen("tcp", address)
	errCh := make(chan error, 1)

	if err != nil {
		return server, nil, err
	}

	go func() {
		defer close(errCh)

		if err := server.Serve(listener); err != nil {
			errCh <- err
		}
	}()

	return server, errCh, nil
}

// runWorkers starts the application's background workers and returns a channel
// through which their terminal error is reported.
func runWorkers(ctx context.Context, app *bootstrap.App) <-chan error {
	errCh := make(chan error, 3)

	go func() {
		errCh <- app.Workers.DownloadTask.Run(ctx)
	}()

	go func() {
		errCh <- app.Workers.InspectTask.Run(ctx)
	}()

	go func() {
		errCh <- app.Workers.OutboxEvent.Run(ctx)
	}()

	return errCh
}

// shutdownHTTP gracefully shuts down the HTTP server within the configured timeout.
func shutdownHTTP(server *http.Server) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
