package router

import (
	"github.com/gorilla/mux"
	"net/http"
	"net/http/pprof"
	"***REMOVED***/darkroom/server/config"
	"***REMOVED***/darkroom/server/handler"
	"***REMOVED***/darkroom/server/service"
)

func NewRouter(deps *service.Dependencies) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	r.Methods(http.MethodGet).Path("/ping").Handler(handler.Ping())

	if config.DebugModeEnabled() {
		setDebugRoutes(r)
	}

	// Catch all handler
	r.Methods(http.MethodGet).PathPrefix("/").Handler(handler.ImageHandler(deps))

	return r
}

func setDebugRoutes(r *mux.Router) {
	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)
	r.HandleFunc("/debug/pprof/goroutine", pprof.Index)
	r.HandleFunc("/debug/pprof/heap", pprof.Index)
	r.HandleFunc("/debug/pprof/threadcreate", pprof.Index)
}
