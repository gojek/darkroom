package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"***REMOVED***/darkroom/server/config"
	"***REMOVED***/darkroom/server/logger"
	"***REMOVED***/darkroom/server/router"
	"***REMOVED***/darkroom/server/service"
	"syscall"
)

func listenServer(s *http.Server) {
	err := s.ListenAndServe()
	if err != http.ErrServerClosed && err != nil {
		logger.Errorf("error while starting %s server: %s", config.AppName(), err)
	}
}

func waitForShutdown(s *http.Server) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig,
		syscall.SIGINT,
		syscall.SIGTERM)
	_ = <-sig
	logger.Infof("%s server shutting down", config.AppName())

	err := s.Shutdown(context.Background())
	if err != nil {
		logger.Error(err.Error())
	}
	logger.Infof("%s server shutdown complete", config.AppName())
}

func Start() {
	logger.Infof("Starting %s server", config.AppName())

	muxRouter := router.NewRouter(service.NewDependencies())

	portInfo := fmt.Sprintf(":%d", config.Port())
	server := &http.Server{Addr: portInfo, Handler: muxRouter}

	go listenServer(server)
	waitForShutdown(server)
}
