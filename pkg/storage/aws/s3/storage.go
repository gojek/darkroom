package s3

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	"github.com/gojek/darkroom/pkg/storage"
)

// Storage holds the fields used by S3 storage implementation
type Storage struct {
	bucketName   string
	bucketRegion string
	accessKey    string
	secretKey    string
	hystrixCmd   storage.HystrixCommand
	downloader   s3manageriface.DownloaderAPI
}

// Get takes in the Context and path as an argument and returns an IResponse interface implementation.
// This method figures out how to get the data from the S3 storage backend.
func (s *Storage) Get(ctx context.Context, path string) storage.IResponse {
	buff := &aws.WriteAtBuffer{}

	responseChannel := make(chan error, 1)
	makeNetworkCall(s.hystrixCmd.Name, s.hystrixCmd.Config, func() error {
		_, err := s.downloader.Download(buff, &s3.GetObjectInput{
			Bucket: aws.String(s.bucketName),
			Key:    aws.String(path),
		})
		responseChannel <- err
		return err
	}, func(e error) error {
		responseChannel <- e
		return e
	})
	s3Err := <-responseChannel

	return storage.NewResponse([]byte(buff.Bytes()), getStatusCodeFromError(s3Err), s3Err)
}

// NewStorage returns a new s3.Storage instance
func NewStorage(opts ...Option) *Storage {
	s := Storage{}
	for _, opt := range opts {
		opt(&s)
	}
	cfg := aws.NewConfig().WithRegion(s.bucketRegion).WithCredentials(
		credentials.NewStaticCredentials(s.accessKey, s.secretKey, ""),
	)
	ssn, _ := session.NewSession(cfg)
	s.downloader = s3manager.NewDownloaderWithClient(s3.New(ssn))
	return &s
}
