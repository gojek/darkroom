package cmd

import (
	"github.com/gojek/darkroom/cmd/signals"
	"github.com/spf13/cobra"
)

// newRootCmd represents the base command when called without any subcommands.
func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "darkroom",
		Short: "Darkroom is an Image Proxy on your image source",
	}
	cmd.AddCommand(newRunCmdWithOpts(runCmdOpts{
		SetupSignalHandler: signals.SetupSignalHandler,
	}))
	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return newRootCmd().Execute()
}
