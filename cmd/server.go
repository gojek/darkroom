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

func startServer() {
	handler := router.NewRouter(service.NewDependencies())
	s := server.NewServer(server.WithHandler(handler))
	s.Start()
}
