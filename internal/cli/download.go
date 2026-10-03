package cli

import "github.com/spf13/cobra"

func (rt *runtime) newDownloadCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "download",
		Short: "Download a file from a URL",
		RunE:  rt.runDownload,
	}
}

func (rt *runtime) runDownload(cmd *cobra.Command, args []string) error {}
