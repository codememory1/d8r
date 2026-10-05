package command

import (
	"context"

	"github.com/codememory1/d8r/internal/bootstrap"
	"github.com/spf13/cobra"
)

type WithApp func(
	cmd *cobra.Command,
	handle func(context.Context, *bootstrap.App) error,
) error
