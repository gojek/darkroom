package s3

import (
	"net/http"
	"regexp"
	"strconv"
)

var re = regexp.MustCompile("status code: ([0-9]{3})")

func getStatusCodeFromError(err error) int {
	if err == nil {
		return http.StatusOK
	}
	if match := re.FindStringSubmatch(err.Error()); match != nil {
		val, _ := strconv.Atoi(match[1])
		return val
	}
	return http.StatusUnprocessableEntity
}
