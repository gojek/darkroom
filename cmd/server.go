package cmd

import (
	"github.com/gojek/darkroom/pkg/router"
	"github.com/gojek/darkroom/pkg/server"
	"github.com/gojek/darkroom/pkg/service"
	"github.com/spf13/cobra"
)

type runCmdOpts struct {
	SetupSignalHandler func() (stopCh <-chan struct{})
}

func newRunCmdWithOpts(opts runCmdOpts) *cobra.Command {
	args := struct {
		port int
	}{}
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start the app server",
		RunE: func(cmd *cobra.Command, _ []string) error {
			deps, err := service.NewDependencies()
			if err != nil {
				return err
			}
			handler := router.NewRouter(deps)
			s := server.NewServer(server.Options{
				Handler: handler,
				Port:    args.port,
			})
			return s.Start(opts.SetupSignalHandler())
		},
	}
	cmd.PersistentFlags().IntVarP(&args.port, "port", "p", 3000, "server port")
	return cmd
}
