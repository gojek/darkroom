package processor

// Processor interface for performing operations on image
type Processor interface {
	Crop(input []byte, width, height int, point CropPoint) ([]byte, error)
	Resize(input []byte, width, height int) ([]byte, error)
	Watermark(base []byte, overlay []byte, opacity uint8) ([]byte, error)
	GrayScale(input []byte) ([]byte, error)
}
