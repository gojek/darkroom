package handler

import (
	"fmt"
	"net/http"
	"***REMOVED***/darkroom/server/config"
	"***REMOVED***/darkroom/server/constants"
	"***REMOVED***/darkroom/server/service"
)

func ImageHandler(deps *service.Dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := deps.Storage.Get(r.Context(), r.URL.Path)
		if res.Error() != nil {
			// TODO Handle error
			return
		}

		cl, _ := w.Write([]byte(res.Data()))
		w.Header().Set(constants.ContentLengthHeader, string(cl))
		w.Header().Set(constants.CacheControlHeader, fmt.Sprintf("public,max-age=%d", config.CacheTime()))
	}
}
