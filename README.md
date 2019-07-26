# Darkroom - Yet Another Image Proxy

<p align="center"><img src="docs/darkroom-logo.png" width="360"></p>

[![build status](https://travis-ci.com/gojek/darkroom.svg?branch=master)](https://travis-ci.com/gojek/darkroom)
[![Coverage Status](https://coveralls.io/repos/github/gojek/darkroom/badge.svg?branch=master)](https://coveralls.io/github/gojek/darkroom?branch=master)
[![GoDoc](https://godoc.org/github.com/gojek/darkroom?status.svg)](https://godoc.org/github.com/gojek/darkroom)
[![Go Report Card](https://goreportcard.com/badge/github.com/gojek/darkroom)](https://goreportcard.com/report/github.com/gojek/darkroom)
[![GolangCI](https://golangci.com/badges/github.com/gojek/darkroom.svg)](https://golangci.com)
[![GitHub Release](https://img.shields.io/github/release/gojek/darkroom.svg?style=flat)](https://github.com/gojek/darkroom/releases)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)  

#### About
Darkroom combines the [storage backend](pkg/storage) and the [image processor](pkg/processor) and acts as an `Image Proxy` on your image source.  
You may implement your own `Storage` and `Processor` interfaces to gain custom functionality while still keeping other Darkroom Server functionality.  
The native implementations focus on speed and resiliency.

#### Installation
```bash
go get -u github.com/gojek/darkroom
```

### Features
Darkroom acts as an image proxy and currently support several image processing operations such as:
- Cropping based on given anchor points (top, left, right, bottom)
- Resizing
- Grayscaling

### Running the Image Proxy Service
The project has docker images available. They can be tested locally or can be be deployed to production.

Create a file containing the environment variables mentioned in [`config.yaml.example`](./config.yaml.example) and save it as `config.env`
> Note: Bucket credentials are dummy, you need to provide your own credentials.
```bash
DEBUG=true
LOG_LEVEL=debug

APP_NAME=darkroom
APP_VERSION=0.0.1
APP_DESCRIPTION=A realtime image processing service

SOURCE_KIND=s3
SOURCE_BUCKET_NAME=bucket-name
SOURCE_BUCKET_REGION=bucket-region
SOURCE_BUCKET_ACCESSKEY=AKIA*************
SOURCE_BUCKET_SECRETKEY=4y/*******************************
SOURCE_PATHPREFIX=/uploads

PORT=3000

CACHE_TIME=31536000

SOURCE_HYSTRIX_COMMANDNAME=S3_ADAPTER
SOURCE_HYSTRIX_TIMEOUT=5000
SOURCE_HYSTRIX_MAXCONCURRENTREQUESTS=100
SOURCE_HYSTRIX_REQUESTVOLUMETHRESHOLD=10
SOURCE_HYSTRIX_SLEEPWINDOW=10
SOURCE_HYSTRIX_ERRORPERCENTTHRESHOLD=25
```
Build the docker image and run it with the config created.
```bash
make docker-image
docker run -p 80:3000 --env-file ./config.env ${USER}/darkroom:latest
```

### Customisation
You may customise darkroom, for example, you may want to write a storage backend that talks to another service and gets the images.
Or might want to create an image processor that uses GPU acceleration to speed up the performance.
#### Available Interfaces
```go
type Processor interface {
	Crop(img image.Image, width, height int, point CropPoint) image.Image
	Decode(data []byte) (image.Image, string, error)
	Encode(img image.Image, format string) ([]byte, error)
	GrayScale(img image.Image) image.Image
	Resize(img image.Image, width, height int) image.Image
	Watermark(base []byte, overlay []byte, opacity uint8) ([]byte, error)
	Flip(image image.Image, mode string) image.Image
	Rotate(image image.Image, angle float64) image.Image
	FixOrientation(img image.Image, orientation int) image.Image
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

