package router

import (
	"log"
	"net/http"
	"net/http/pprof"

	"***REMOVED***/darkroom/core/pkg/regex"

	"github.com/gorilla/mux"
	"***REMOVED***/darkroom/core/internal/handler"
	"***REMOVED***/darkroom/core/pkg/config"
	"***REMOVED***/darkroom/core/service"
)

func NewRouter(deps *service.Dependencies) *mux.Router {
	validateDependencies(deps)
	r := mux.NewRouter().StrictSlash(true)

	r.Methods(http.MethodGet).Path("/ping").Handler(handler.Ping())

	if config.DebugModeEnabled() {
		setDebugRoutes(r)
	}

	// Catch all handler
	s := config.Source()
	if (regex.S3Matcher.MatchString(s.Kind) ||
		regex.CloudfrontMatcher.MatchString(s.Kind)) &&
		s.PathPrefix != "" {
		r.Methods(http.MethodGet).PathPrefix(s.PathPrefix).Handler(handler.ImageHandler(deps))
	} else {
		r.Methods(http.MethodGet).PathPrefix("/").Handler(handler.ImageHandler(deps))
	}

	return r
}

func validateDependencies(deps *service.Dependencies) {
	if deps.Storage == nil || deps.Manipulator == nil {
		log.Fatal("handler dependencies are not valid")
	}
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
