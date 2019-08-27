---
id: customization
title: Customization
---
## Customization
You may customize darkroom, for example, you may want to write a storage backend that talks to another service and gets the images.
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

## Custom Storage Example

This example shows how you can implement a custom `Storage` which reads file from the local storage of the machine on which darkroom server is running.

- Create file `pkg/storage/local/storage.go`
```go
package local

// Storage holds the fields used by local storage implementation
type Storage struct {
	volume string
}

// Option represents the local storage options
type Option func(s *Storage)

// WithVolume sets the volume
func WithVolume(volume string) Option {
	return func(s *Storage) {
		s.volume = volume
	}
}

// Get takes in the Context and path as an argument and returns an IResponse interface implementation.
// This method figures out how to get the data from the local storage backend.
func (s *Storage) Get(ctx context.Context, path string) storage.IResponse {
	data, err := ioutil.ReadFile(fmt.Sprintf("%s%s", s.volume, path))
	if err != nil {
		return storage.NewResponse([]byte(err.Error()), http.StatusUnprocessableEntity, err)
	}
	return storage.NewResponse([]byte(data), http.StatusOK, nil)
}

// NewStorage returns a new local.Storage instance
func NewStorage(opts ...Option) *Storage {
	s := Storage{}
	for _, opt := range opts {
		opt(&s)
	}
	return &s
}
```


- Inject this `local.Storage` in the handler dependencies `pkg/service/dependencies.go`
```go
func NewDependencies() *Dependencies {
	deps := &Dependencies{Manipulator: NewManipulator(native.NewBildProcessor())}
	deps.Storage = local.NewStorage(
		local.WithVolume("/absolute/path/to/folder/containing/images"),
	)
	return deps
}
```


When you make a call to `http://localhost:3000/sample-image.jpg?w=500`, Darkroom will try to get the image from the location `/absolute/path/to/folder/containing/images/sample-image.jpg` on the local disk and serve a 500 pixel width image.
