package config

import (
	"***REMOVED***/darkroom/core/pkg/regex"
	"***REMOVED***/darkroom/core/pkg/storage"
)

type app struct {
	name        string
	version     string
	description string
}

type S3Bucket struct {
	Name      string
	Region    string
	AccessKey string
	SecretKey string
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
	PathPrefix     string
}

func (s *source) readValue() {
	v := Viper()
	if regex.S3Matcher.MatchString(s.Kind) {
		s.Value = S3Bucket{
			Name:      v.GetString("source.bucket.name"),
			Region:    v.GetString("source.bucket.region"),
			AccessKey: v.GetString("source.bucket.accessKey"),
			SecretKey: v.GetString("source.bucket.secretKey"),
		}
	} else if regex.CloudfrontMatcher.MatchString(s.Kind) {
		s.Value = Cloudfront{
			Host:           v.GetString("source.host"),
			SecureProtocol: v.GetBool("source.secureProtocol"),
		}
	} else if regex.WebFolderMatcher.MatchString(s.Kind) {
		s.Value = WebFolder{BaseURL: v.GetString("source.baseURL")}
	}
}
