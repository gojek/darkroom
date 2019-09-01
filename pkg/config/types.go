package config

import (
	"github.com/gojek/darkroom/pkg/regex"
	"github.com/gojek/darkroom/pkg/storage"
)

// S3Bucket contains the configuration values for S3 source
type S3Bucket struct {
	// Name of the bucket
	Name string
	// Region of the bucket
	Region string
	// Access key that should be used to access the bucket
	AccessKey string
	// Secret key that should be used to access the bucket
	SecretKey string
}

// WebFolder contains the configuration for a directory available on the public internet
type WebFolder struct {
	// BaseURL that should be appended to the path
	// Eg: https://example.com/web-folder/{path} will map to https://host.com/{path}
	BaseURL string
}

// Cloudfront contains the configuration for cloudfront which can be used with an S3 bucket
type Cloudfront struct {
	// Host is the FQDN for the cloudfront integration on the S3 bucket
	Host string
	// SecureProtocol designates whether to use http or https protocol for requests
	SecureProtocol bool
}

// Source contains the configuration for data source object that will be used, the type of the data source, hystrix command,
// and the path prefix to restring access
type Source struct {
	// Kind tells which type of Storage backend should be used
	Kind string
	// HystrixCommand provides the hystrix config to be used with the source to add resiliency
	HystrixCommand storage.HystrixCommand
	// Value is and interface which holds the actual kind of the object
	Value interface{}
	// PathPrefix is used to restrict access to specific paths only via the image proxy
	PathPrefix string
}

func (s *Source) readValue() {
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
