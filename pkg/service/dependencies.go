package service

import (
	"time"

	"github.com/gojektech/heimdall"
	"github.com/gojektech/heimdall/hystrix"
	"github.com/gojek/darkroom/pkg/config"
	"github.com/gojek/darkroom/pkg/processor/native"
	"github.com/gojek/darkroom/pkg/regex"
	base "github.com/gojek/darkroom/pkg/storage"
	"github.com/gojek/darkroom/pkg/storage/aws/cloudfront"
	"github.com/gojek/darkroom/pkg/storage/aws/s3"
	"github.com/gojek/darkroom/pkg/storage/webfolder"
)

type Dependencies struct {
	Storage     base.Storage
	Manipulator Manipulator
}

func NewDependencies() *Dependencies {
	s := config.Source()
	deps := &Dependencies{Manipulator: NewManipulator(native.NewBildProcessor())}
	if regex.WebFolderMatcher.MatchString(s.Kind) {
		deps.Storage = NewWebFolderStorage(s.Value.(config.WebFolder), s.HystrixCommand)
	} else if regex.S3Matcher.MatchString(s.Kind) {
		deps.Storage = NewS3Storage(s.Value.(config.S3Bucket), s.HystrixCommand)
	} else if regex.CloudfrontMatcher.MatchString(s.Kind) {
		deps.Storage = NewCloudfrontStorage(s.Value.(config.Cloudfront), s.HystrixCommand)
	}
	return deps
}

func NewS3Storage(b config.S3Bucket, hc base.HystrixCommand) *s3.Storage {
	return s3.NewStorage(
		s3.WithBucketName(b.Name),
		s3.WithBucketRegion(b.Region),
		s3.WithAccessKey(b.AccessKey),
		s3.WithSecretKey(b.SecretKey),
		s3.WithHystrixCommand(hc),
	)
}

func NewWebFolderStorage(wf config.WebFolder, hc base.HystrixCommand) *webfolder.Storage {
	return webfolder.NewStorage(
		webfolder.WithBaseURL(wf.BaseURL),
		webfolder.WithHeimdallClient(newHystrixClient(hc)),
	)
}

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
