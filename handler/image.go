package handler

import (
	"fmt"
	"net/http"

	"***REMOVED***/darkroom/core/config"
	"***REMOVED***/darkroom/core/logger"
	"***REMOVED***/darkroom/core/service"
)

const (
	ContentLengthHeader = "Content-Length"
	CacheControlHeader  = "Cache-Control"
)

func ImageHandler(deps *service.Dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := deps.Storage.Get(r.Context(), r.URL.Path)
		if res.Error() != nil {
			logger.Errorf("error: %s", res.Error())
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
				logger.Errorf("error: %s", res.Error())
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
		}

		cl, _ := w.Write([]byte(data))
		w.Header().Set(ContentLengthHeader, string(cl))
		w.Header().Set(CacheControlHeader, fmt.Sprintf("public,max-age=%d", config.CacheTime()))
	}
}
