package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "darkroom",
	Short: "Darkroom is an Image Proxy on your image source",
}

func init() {}

// Run function lets you run the commands
func Run(args []string) error {
	rootCmd.SetArgs(args)
	return rootCmd.Execute()
}
