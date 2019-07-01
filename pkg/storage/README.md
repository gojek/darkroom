# Storage Backend for Darkroom

#### About
This module holds the logic to get an image blob from a source. It is used by the Darkroom [Application Server](https://github.com/gojek/darkroom).  
You may implement the `Storage` interface to gain custom functionality while still keeping other Darkroom functionality.  
You may write custom backend for downloading images from a hosting provider or a web proxy.

#### Interfaces
```go
type Storage interface {
	Get(ctx context.Context, path string) IResponse
}
```
Any `struct` implementing the above interface can be used with Darkroom.
> Note: The response returned should be of type `IResponse`
```go
type IResponse interface {
	Data() []byte
	Error() error
	Status() int
}
```

#### Example

```go
s3s := s3.NewStorage(
                s3.WithBucketName("bucket"),
                s3.WithBucketRegion("region"),
                s3.WithAccessKey("randomAccessKey"),
                s3.WithSecretKey("randomSecretKey"),
                s3.WithHystrixCommand(storage.HystrixCommand{
                    Name: "TestCommand",
                    Config: hystrix.CommandConfig{
                        Timeout:                2000,
                        MaxConcurrentRequests:  100,
                        RequestVolumeThreshold: 10,
                        SleepWindow:            10,
                        ErrorPercentThreshold:  25,
                    },
                }),
            )

res := s3s.Get(context.Background(), "path/to/image-blob")
if res.Error() != nil {
	fatal(res.Error())
}

// res.Status() -> Status Code returned by the backend
// res.Data() -> Data in form of []byte returned by the backend
```
