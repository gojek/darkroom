package main

import (
	"github.com/urfave/cli"
	"os"
	"***REMOVED***/darkroom/server/config"
	"***REMOVED***/darkroom/server/logger"
	"***REMOVED***/darkroom/server/router"
	"***REMOVED***/darkroom/server/server"
	"***REMOVED***/darkroom/server/service"
)

func main() {
	defer func() {
		if e := recover(); e != nil {
			logger.Errorf("failed to start the app due to error: %s", e)
		}
	}()

	a := cli.NewApp()
	a.Name = config.AppName()
	a.Version = config.AppVersion()
	a.Usage = config.AppDescription()
	a.Action = func(c *cli.Context) error {
		handler := router.NewRouter(service.NewDependencies())
		s := server.NewServer(server.WithHandler(handler))
		s.Start()
		return nil
	}

	if err := a.Run(os.Args); err != nil {
		panic(err)
	}
}
