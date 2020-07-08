package cmd

import (
	"github.com/gojek/darkroom/internal/version"
	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version",
		Long:  `Print version.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			buildInfo := version.Build
			cmd.Println(buildInfo.Version)
			return nil
		},
	}
}
