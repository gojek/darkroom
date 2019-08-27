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
	// Blur takes an input byte array and returns the blurred byte array by the specified
	// radius(<=1000) or error radius must be larger than 0
	Blur(image image.Image, radius float64) image.Image
	// Watermark takes an input byte array, overlay byte array and opacity value
	// and returns the watermarked image bytes or error
	Watermark(base []byte, overlay []byte, opacity uint8) ([]byte, error)
	// Flip takes an input image and returns the image flipped. The direction of flip
	// is determined by the specified mode - 'v' for a vertical flip, 'h' for a horizontal flip and
	// 'vh'(or 'hv') for both.
	Flip(image image.Image, mode string) image.Image
	// Rotate takes an input image and returns a image rotated by the specified degrees.
	// The rotation is applied clockwise, and fractional angles are supported.
	Rotate(image image.Image, angle float64) image.Image
	// Decode takes a byte array and returns the image, extension, and error
	Decode(data []byte) (img image.Image, format string, err error)
	// Encode takes an image and extension and return the encoded byte array or error
	Encode(img image.Image, format string) ([]byte, error)
	// FixOrientation takes an image and it's EXIF orientation (if exist)
	// and returns the image with its EXIF orientation fixed
	FixOrientation(img image.Image, orientation int) image.Image
}
