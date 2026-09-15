package cli

import (
	"context"
	"log/slog"

	"github.com/codememory1/d8r/internal/bootstrap"
	"github.com/codememory1/d8r/internal/infrastructure/config"
	"github.com/spf13/cobra"
)

var (
	// Default path to the application configuration
	defaultConfiguration = "./config/config.yaml"
)

type runtime struct {
	configuration string
}

// Execute creates and configures the application's root CLI command,
// registers global flags, and executes it with the provided context.
func Execute(ctx context.Context) error {
	rt := &runtime{}
	root := &cobra.Command{
		Use:           "d8r",
		Short:         "Operate the commerce modular monolith",
		Version:       "1.0.0",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	
	root.AddCommand(rt.newServerCommand())

	root.PersistentFlags().StringVarP(
		&rt.configuration,
		"configuration",
		"c",
		defaultConfiguration,
		"configuration YAML path",
	)

	return root.ExecuteContext(ctx)
}

// withApp loads the configuration, initializes the logger and the App,
// and then passes the ready App instance to the handler.
func (rt *runtime) withApp(cmd *cobra.Command, handle func(ctx context.Context, app *bootstrap.App) error) error {
	configuration, err := config.Load(rt.configuration)

	if err != nil {
		return err
	}

	logger := slog.New(slog.NewJSONHandler(cmd.ErrOrStderr(), &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	app, err := bootstrap.NewApp(cmd.Context(), configuration, logger)

	if err != nil {
		return err
	}

	return handle(cmd.Context(), app)
}
