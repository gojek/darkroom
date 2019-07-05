package processor

import "image"

// Processor interface for performing operations on image bytes
type Processor interface {
	// Crop takes an image.Image, width, height and a CropPoint and returns the cropped image
	Crop(image image.Image, width, height int, point CropPoint) image.Image
	// Resize takes an image.Image, width and height and returns the re-sized image
	Resize(image image.Image, width, height int) image.Image
	// GrayScale takes an input byte array and returns the grayscaled byte array or error
	GrayScale(image image.Image) image.Image
	// Watermark takes an input byte array, overlay byte array and opacity value
	// and returns the watermarked image bytes or error
	Watermark(base []byte, overlay []byte, opacity uint8) ([]byte, error)
	// Decode takes a byte array and returns the image, extension, and error
	Decode(data []byte) (image.Image, string, error)
	// Encode takes an image and extension and return the encoded byte array or error
	Encode(img image.Image, format string) ([]byte, error)
}
