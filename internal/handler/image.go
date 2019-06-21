package handler

import (
	"fmt"
	"net/http"
	"***REMOVED***/darkroom/core/pkg/metrics"

	"***REMOVED***/darkroom/core/pkg/config"
	"***REMOVED***/darkroom/core/pkg/logger"
	"***REMOVED***/darkroom/core/pkg/service"
)

const (
	ContentLengthHeader = "Content-Length"
	CacheControlHeader  = "Cache-Control"
	StorageGetErrorKey  = "handler.storage.get.error"
	ProcessorErrorKey   = "handler.processor.error"
)

func ImageHandler(deps *service.Dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.SugaredWithRequest(r)
		res := deps.Storage.Get(r.Context(), r.URL.Path)
		if res.Error() != nil {
			l.Errorf("error from Storage.Get: %s", res.Error())
			metrics.Update(metrics.UpdateOption{Name: StorageGetErrorKey, Type: metrics.Count})
			w.WriteHeader(res.Status())
			return
		}
		var data []byte
		var err error
		data = res.Data()

		params := make(map[string]string)
		values := r.URL.Query()
		if len(values) > 0 {
			for v := range values {
				if len(values.Get(v)) != 0 {
					params[v] = values.Get(v)
				}
			}
			data, err = deps.Manipulator.Process(r.Context(), data, params)
			if err != nil {
				l.Errorf("error from Manipulator.Process: %s", err)
				metrics.Update(metrics.UpdateOption{Name: ProcessorErrorKey, Type: metrics.Count})
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
		}

		cl, _ := w.Write([]byte(data))
		w.Header().Set(ContentLengthHeader, string(cl))
		w.Header().Set(CacheControlHeader, fmt.Sprintf("public,max-age=%d", config.CacheTime()))
	}
}
