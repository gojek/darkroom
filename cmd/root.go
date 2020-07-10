package cmd

import (
	"github.com/spf13/cobra"
)

// newRootCmd represents the base command when called without any subcommands.
func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "darkroom",
		Short: "Darkroom is an Image Proxy on your image source",
	}
	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return newRootCmd().Execute()
}
