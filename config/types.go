package config

import (
	"***REMOVED***/darkroom/core/constants"
	"***REMOVED***/darkroom/storage"
)

type loggerConfig struct {
	level  string
	format string
}

type app struct {
	name        string
	version     string
	description string
}

type S3Bucket struct {
	Name       string
	Region     string
	AccessKey  string
	SecretKey  string
	PathPrefix string
}

type WebFolder struct {
	BaseURL string
}

type Cloudfront struct {
	Host           string
	SecureProtocol bool
}

type source struct {
	Kind           string
	HystrixCommand storage.HystrixCommand
	Value          interface{}
}

func (s *source) readValue() {
	v := Viper()
	if constants.S3Matcher.MatchString(s.Kind) {
		s.Value = S3Bucket{
			Name:       v.GetString("source.bucket.name"),
			Region:     v.GetString("source.bucket.region"),
			AccessKey:  v.GetString("source.bucket.accessKey"),
			SecretKey:  v.GetString("source.bucket.secretKey"),
			PathPrefix: v.GetString("source.bucket.pathPrefix"),
		}
	} else if constants.CloudfrontMatcher.MatchString(s.Kind) {
		s.Value = Cloudfront{
			Host:           v.GetString("source.host"),
			SecureProtocol: v.GetBool("source.secureProtocol"),
		}
	} else if constants.WebFolderMatcher.MatchString(s.Kind) {
		s.Value = WebFolder{BaseURL: v.GetString("source.baseURL")}
	}
}
