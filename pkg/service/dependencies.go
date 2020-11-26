// Package service contains the service definitions used by the handler
package service

import (
	"errors"
	"strings"
	"time"

	"github.com/gojektech/heimdall"
	"github.com/gojektech/heimdall/hystrix"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/gojek/darkroom/pkg/config"
	"github.com/gojek/darkroom/pkg/logger"
	"github.com/gojek/darkroom/pkg/metrics"
	"github.com/gojek/darkroom/pkg/processor/native"
	"github.com/gojek/darkroom/pkg/regex"
	base "github.com/gojek/darkroom/pkg/storage"
	"github.com/gojek/darkroom/pkg/storage/aws/cloudfront"
	"github.com/gojek/darkroom/pkg/storage/aws/s3"
	"github.com/gojek/darkroom/pkg/storage/gcs"
	"github.com/gojek/darkroom/pkg/storage/webfolder"
)

// Dependencies struct holds the reference to the Storage and the Manipulator interface implementations
type Dependencies struct {
	Storage       base.Storage
	Manipulator   Manipulator
	MetricService metrics.MetricService
}

// NewDependencies constructs new Dependencies based on the config.DataSource().Kind
// Currently, it supports only one Manipulator
func NewDependencies(registry *prometheus.Registry) (deps *Dependencies, err error) {
	var metricService metrics.MetricService
	if regex.PrometheusMatcher.MatchString(config.MetricsSystem()) {
		metricService = metrics.NewPrometheus(registry)
	} else if regex.StatsdMatcher.MatchString(config.MetricsSystem()) {
		metricService, _ = metrics.InitializeStatsdCollector(config.StatsdConfig())
	}
	if metricService == nil {
		metricService = metrics.NoOpMetricService{}
		logger.Warn("NoOpMetricService is being used since metric system is not specified")
	}
	deps = &Dependencies{
		Manipulator:   NewManipulator(native.NewBildProcessor(), getDefaultParams(), metricService),
		MetricService: metricService,
	}
	s := config.DataSource()
	if regex.WebFolderMatcher.MatchString(s.Kind) {
		deps.Storage = NewWebFolderStorage(s.Value.(config.WebFolder), s.HystrixCommand)
	} else if regex.S3Matcher.MatchString(s.Kind) {
		deps.Storage = NewS3Storage(s.Value.(config.S3Bucket), s.HystrixCommand)
	} else if regex.CloudfrontMatcher.MatchString(s.Kind) {
		deps.Storage = NewCloudfrontStorage(s.Value.(config.Cloudfront), s.HystrixCommand)
	} else if regex.GoogleCloudStorageMatcher.MatchString(s.Kind) {
		deps.Storage, err = NewGoogleCloudStorage(s.Value.(config.GoogleCloudStorage), s.HystrixCommand)
	}
	if deps.Storage == nil || deps.Manipulator == nil {
		return nil, errors.New("handler dependencies are not valid")
	}
	return deps, err
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

// NewGoogleCloudStorage create a new gcs.Storage struct from the config.GoogleCloudStorage and the HystrixCommand
func NewGoogleCloudStorage(b config.GoogleCloudStorage, hc base.HystrixCommand) (*gcs.Storage, error) {
	return gcs.NewStorage(gcs.Options{
		BucketName:      b.Name,
		CredentialsJSON: []byte(b.CredentialsJSON),
		Client:          newHystrixClient(hc),
	})
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
