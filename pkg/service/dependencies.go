// Package service contains the service definitions used by the handler
package service

import (
	"errors"
	"github.com/gojek/darkroom/pkg/metrics"
	"github.com/gojek/darkroom/pkg/regex"
	"github.com/prometheus/client_golang/prometheus"
	"strings"
	"time"

	"github.com/gojek/darkroom/pkg/config"
	"github.com/gojek/darkroom/pkg/processor/native"
	base "github.com/gojek/darkroom/pkg/storage"
	"github.com/gojek/darkroom/pkg/storage/aws/cloudfront"
	"github.com/gojek/darkroom/pkg/storage/aws/s3"
	"github.com/gojek/darkroom/pkg/storage/webfolder"
	"github.com/gojektech/heimdall"
	"github.com/gojektech/heimdall/hystrix"
)

// Dependencies struct holds the reference to the Storage and the Manipulator interface implementations
type Dependencies struct {
	Storage     base.Storage
	Manipulator Manipulator
	Registry *prometheus.Registry
	MetricService metrics.MetricService
}

// NewDependencies constructs new Dependencies based on the config.DataSource().Kind
// Currently, it supports only one Manipulator
func NewDependencies() (*Dependencies, error) {
	s := config.DataSource()
	registry := prometheus.NewRegistry()
	metricService := metrics.NewPrometheus(registry)
	deps := &Dependencies{Manipulator: NewManipulator(native.NewBildProcessor(), getDefaultParams(), metricService),
		MetricService: metricService, Registry: registry}
	if regex.WebFolderMatcher.MatchString(s.Kind) {
		deps.Storage = NewWebFolderStorage(s.Value.(config.WebFolder), s.HystrixCommand)
	} else if regex.S3Matcher.MatchString(s.Kind) {
		deps.Storage = NewS3Storage(s.Value.(config.S3Bucket), s.HystrixCommand)
	} else if regex.CloudfrontMatcher.MatchString(s.Kind) {
		deps.Storage = NewCloudfrontStorage(s.Value.(config.Cloudfront), s.HystrixCommand)
	}
	if deps.Storage == nil || deps.Manipulator == nil {
		return nil, errors.New("handler dependencies are not valid")
	}
	return deps, nil
}

func getDefaultParams() map[string]string {
	params := make(map[string]string)
	for _, param := range config.DefaultParams() {
		if strings.Contains(param, "=") {
			p := strings.Split(param, "=")
			params[p[0]] = p[1]
		}
	}
	return params
}

// NewS3Storage create a new s3.Storage struct from the config.S3Bucket and the HystrixCommand
func NewS3Storage(b config.S3Bucket, hc base.HystrixCommand) *s3.Storage {
	return s3.NewStorage(
		s3.WithBucketName(b.Name),
		s3.WithBucketRegion(b.Region),
		s3.WithAccessKey(b.AccessKey),
		s3.WithSecretKey(b.SecretKey),
		s3.WithHystrixCommand(hc),
	)
}

// NewWebFolderStorage create a new webfolder.Storage struct from the config.WebFolder and the HystrixCommand
func NewWebFolderStorage(wf config.WebFolder, hc base.HystrixCommand) *webfolder.Storage {
	return webfolder.NewStorage(
		webfolder.WithBaseURL(wf.BaseURL),
		webfolder.WithHeimdallClient(newHystrixClient(hc)),
	)
}

// NewCloudfrontStorage create a new cloudfront.Storage struct from the config.Cloudfront and the HystrixCommand
func NewCloudfrontStorage(c config.Cloudfront, hc base.HystrixCommand) *cloudfront.Storage {
	var opts []cloudfront.Option
	if c.SecureProtocol {
		opts = append(opts, cloudfront.WithSecureProtocol())
	}
	opts = append(opts,
		cloudfront.WithCloudfrontHost(c.Host),
		cloudfront.WithHeimdallClient(newHystrixClient(hc)),
	)
	return cloudfront.NewStorage(opts...)
}

func newHystrixClient(hc base.HystrixCommand) heimdall.Client {
	return hystrix.NewClient(
		hystrix.WithHTTPTimeout(time.Duration(hc.Config.Timeout)*time.Millisecond),
		hystrix.WithMaxConcurrentRequests(hc.Config.MaxConcurrentRequests),
		hystrix.WithRequestVolumeThreshold(hc.Config.RequestVolumeThreshold),
		hystrix.WithSleepWindow(hc.Config.SleepWindow),
		hystrix.WithErrorPercentThreshold(hc.Config.ErrorPercentThreshold),
	)
}
