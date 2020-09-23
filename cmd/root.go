package cmd

import (
	"github.com/gojek/darkroom/internal/runtime"
	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

// newRootCmd represents the base command when called without any subcommands.
func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "darkroom",
		Short: "Darkroom is an Image Proxy on your image source",
	}
	cmd.AddCommand(newRunCmdWithOpts(runCmdOpts{
		SetupSignalHandler: signals.SetupSignalHandler,
		registry:           runtime.PrometheusRegistry(),
	}))
	cmd.AddCommand(newVersionCmd())
	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return newRootCmd().Execute()
}
