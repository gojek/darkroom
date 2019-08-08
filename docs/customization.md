---
id: customization
title: Customization
---
## Customisation
You may customise darkroom, for example, you may want to write a storage backend that talks to another service and gets the images.
Or might want to create an image processor that uses GPU acceleration to speed up the performance.
## Available Interfaces
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
	FixOrientation(image image.Image, orientation int) image.Image
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

