# Application Server for Darkroom
[![build status](https://***REMOVED***/darkroom/core/badges/master/build.svg)](https://***REMOVED***/darkroom/core/commits/master)
[![coverage report](https://***REMOVED***/darkroom/core/badges/master/coverage.svg)](https://***REMOVED***/darkroom/core/commits/master)

#### About
This project combines the darkroom [storage backend](https://***REMOVED***/darkroom/storage) and the [image processor](https://***REMOVED***/darkroom/processor) and acts as an `Image Proxy` on your image source.  
You may implement your own `Storage` and `Processor` interfaces to gain custom functionality while still keeping other Darkroom Server functionality.

#### Installation
The project has docker images available. They can be tested locally or can be be deployed to production.

Create a file containing the environment variables mentioned in [`config.yaml.example`](./config.yaml.example) and save it as `config.env`
> Note: Bucket credentials are dummy, you need to provide your own credentials.
```bash
DEBUG=true
LOG_LEVEL=debug
LOG_FORMAT=json

APP_NAME=darkroom
APP_VERSION=0.0.1
APP_DESCRIPTION=A realtime image processing service

BUCKET_NAME=bucket-name
BUCKET_REGION=bucket-region
BUCKET_ACCESSKEY=AKIAIGTIBDFJ4UIBJDYXQ
BUCKET_SECRETKEY=4y/caOkhg324kgk342hkh3w4/iHLKJGHhkjl4hjkhKG
BUCKET_PATHPREFIX=/base-folder

PORT=3000

CACHE_TIME=31536000

HYSTRIX_COMMAND_NAME=S3_ADAPTER
HYSTRIX_CONFIG_TIMEOUT=5000
HYSTRIX_CONFIG_MAXCONCURRENTREQUESTS=100
HYSTRIX_CONFIG_REQUESTVOLUMETHRESHOLD=10
HYSTRIX_CONFIG_SLEEPWINDOW=10
HYSTRIX_CONFIG_ERRORPERCENTTHRESHOLD=25
```
Build the docker image and run it with the config created.
```bash
docker build --build-arg SSH_PRIVATE_KEY="$(cat ~/.ssh/id_rsa)" -t ${USER}/darkroom:latest .
docker run -p 80:3000 --env-file ./config.env ${USER}/darkroom:latest
```

### Customisation
You may customise darkroom, for example, you may want to write a storage backend that talks to another service and gets the images.
Or might want to create an image processor that uses GPU acceleration to speed up the performance.
#### Available Interfaces
```go
type Processor interface {
	Crop(input []byte, width, height int, point CropPoint) ([]byte, error)
	Resize(input []byte, width, height int) ([]byte, error)
	Watermark(base []byte, overlay []byte, opacity uint8) ([]byte, error)
	GrayScale(input []byte) ([]byte, error)
}
```
```go
type Storage interface {
	Get(ctx context.Context, path string) IResponse
}
```
```go
type IResponse interface {
	Data() []byte
	Error() error
	Status() int
}
```
Any `struct` implementing the above interfaces can be used with Darkroom.  
> Note: The struct implementing the `Storage` interface must return a struct implementing the `IResponse` interface.

