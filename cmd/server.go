package cmd

import (
	"github.com/gojek/darkroom/pkg/logger"
	"github.com/gojek/darkroom/pkg/router"
	"github.com/gojek/darkroom/pkg/server"
	"github.com/gojek/darkroom/pkg/service"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the app server",
	Run:   serverCmdF,
}

func serverCmdF(cmd *cobra.Command, args []string) {
	defer func() {
		if e := recover(); e != nil {
			logger.Errorf("failed to start the app due to error: %s", e)
		}
	}()
	startServer()
}

type runCmdOpts struct {
	SetupSignalHandler func() (stopCh <-chan struct{})
}

func newRunCmdWithOpts(opts runCmdOpts) *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Start the app server",
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := service.NewDependencies()
			if err != nil {
				return err
			}
			handler := router.NewRouter(deps)
			s := server.NewServer(server.WithHandler(handler))
			s.Start()
			return nil
		},
	}
}

func startServer() {
	deps, _ := service.NewDependencies()
	handler := router.NewRouter(deps)
	s := server.NewServer(server.WithHandler(handler))
	s.Start()
}
