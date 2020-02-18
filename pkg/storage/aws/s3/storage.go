package s3

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
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
	service      s3iface.S3API
	hystrixCmd   storage.HystrixCommand
	downloader   s3manageriface.DownloaderAPI
}

// Get takes in the Context and path as an argument and returns an IResponse interface implementation.
// This method figures out how to get the data from the S3 storage backend.
func (s *Storage) Get(ctx context.Context, path string, opt *storage.GetRequestOptions) storage.IResponse {
	input := s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(path),
	}

	if opt != nil && len(opt.Range) > 0 {
		input.Range = &opt.Range
		return s.getRange(ctx, input)
	}

	return s.get(ctx, input)
}

func (s *Storage) get(ctx context.Context, input s3.GetObjectInput) storage.IResponse {
	buff := &aws.WriteAtBuffer{}
	responseChannel := make(chan error, 1)
	makeNetworkCall(s.hystrixCmd.Name, s.hystrixCmd.Config, func() error {
		_, err := s.downloader.Download(buff, &input)
		responseChannel <- err
		return err
	}, func(e error) error {
		responseChannel <- e
		return e
	})
	s3Err := <-responseChannel

	return storage.NewResponse(buff.Bytes(), getStatusCodeFromError(s3Err), s3Err, nil)
}

func (s *Storage) getRange(ctx context.Context, input s3.GetObjectInput) storage.IResponse {
	responseChannel := make(chan *s3.GetObjectOutput, 1)
	errorChannel := make(chan error, 1)
	makeNetworkCall(s.hystrixCmd.Name, s.hystrixCmd.Config, func() error {
		resp, err := s.service.GetObject(&input)
		errorChannel <- err
		responseChannel <- resp
		return err
	}, func(e error) error {
		errorChannel <- e
		return e
	})

	s3Resp, s3Err := <-responseChannel, <-errorChannel

	var metadata *storage.ResponseMetadata
	var body []byte
	if s3Err == nil {
		metadata = s.newMetadata(*s3Resp)
		body, _ = ioutil.ReadAll(s3Resp.Body)
	}
	return storage.NewResponse(body, getStatusCodeFromError(s3Err), s3Err, metadata)
}

func (s *Storage) newMetadata(output s3.GetObjectOutput) *storage.ResponseMetadata {
	metadata := storage.ResponseMetadata{
		AcceptRanges:  aws.StringValue(output.AcceptRanges),
		ContentLength: fmt.Sprintf("%d", aws.Int64Value(output.ContentLength)),
		ContentRange:  aws.StringValue(output.ContentRange),
		ContentType:   aws.StringValue(output.ContentType),
		ETag:          aws.StringValue(output.ETag),
	}

	if output.LastModified != nil {
		metadata.LastModified = aws.TimeValue(output.LastModified).Format(http.TimeFormat)
	}
	return &metadata
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
	s.service = s3.New(ssn)
	s.downloader = s3manager.NewDownloaderWithClient(s.service)
	return &s
}
