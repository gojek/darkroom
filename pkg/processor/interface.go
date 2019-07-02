package processor

// Processor interface for performing operations on image bytes
type Processor interface {
	// Crop takes an input byte array, width, height and a CropPoint and returns the cropped image bytes or error
	Crop(input []byte, width, height int, point CropPoint) ([]byte, error)
	// Resize takes an input byte array, width and height and returns the re-sized image bytes or error
	Resize(input []byte, width, height int) ([]byte, error)
	// Watermark takes an input byte array, overlay byte array and opacity value
	// and returns the watermarked image bytes or error
	Watermark(base []byte, overlay []byte, opacity uint8) ([]byte, error)
	// GrayScale takes an input byte array and returns the grayscaled byte array or error
	GrayScale(input []byte) ([]byte, error)
}
