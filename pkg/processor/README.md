# Image Processor for Darkroom

#### About
This module holds the logic to process images. It is used by the Darkroom [Application Server](https://***REMOVED***/darkroom/core).  
You may implement the `Processor` interface to gain custom functionality while still keeping other Darkroom functionality.

#### Interface
```go
type Processor interface {
	Crop(input []byte, width, height int, point CropPoint) ([]byte, error)
	Resize(input []byte, width, height int) ([]byte, error)
	Watermark(base []byte, overlay []byte, opacity uint8) ([]byte, error)
	GrayScale(input []byte) ([]byte, error)
}
```
Any `struct` implementing the above interface can be used with Darkroom.

#### Example

```go
bp := NewBildProcessor()
img, _ := ioutil.ReadFile("test.png")
output, err := bp.Resize(img, 500, 500)
_, _ := ioutil.WriteFile("output.png", output, 0644)
```