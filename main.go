package main

import (
	"github.com/urfave/cli"
	"os"
	"***REMOVED***/darkroom/server/config"
	"***REMOVED***/darkroom/server/logger"
	"***REMOVED***/darkroom/server/server"
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
		server.Start()
		return nil
	}

	if err := a.Run(os.Args); err != nil {
		panic(err)
	}
}
