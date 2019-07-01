package main

import (
	"os"

	"github.com/urfave/cli"
	"github.com/gojek/darkroom/pkg/config"
	"github.com/gojek/darkroom/pkg/logger"
	"github.com/gojek/darkroom/pkg/router"
	"github.com/gojek/darkroom/pkg/server"
	"github.com/gojek/darkroom/pkg/service"
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
		logger.Errorf("got an error while running main with error: %s", err)
		panic(err)
	}
}
