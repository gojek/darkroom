package native

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"strings"

	"github.com/anthonynsimon/bild/blur"
	"github.com/anthonynsimon/bild/clone"
	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/transform"
	"github.com/gojek/darkroom/pkg/processor"
)

var resizeBoundOption = &transform.RotationOptions{
	ResizeBounds: true,
}

// BildProcessor uses bild library to process images using native Golang image.Image interface
type BildProcessor struct {
	encoders *Encoders
}

// ProcessorOption represents builder function for BildProcessor
type ProcessorOption func(*BildProcessor)

// Crop takes an input image, width, height and a Point and returns the cropped image
func (bp *BildProcessor) Crop(img image.Image, width, height int, point processor.Point) image.Image {
	if width == 0 || height == 0 {
		if width == 0 && height == 0 {
			return img
		}
		return bp.Resize(img, width, height)
	}

	w, h := getResizeWidthAndHeightForCrop(width, height, img.Bounds().Dx(), img.Bounds().Dy())
	img = transform.Resize(img, w, h, transform.Linear)
	x0, y0 := getStartingPointForCrop(w, h, width, height, point)
	rect := image.Rect(x0, y0, width+x0, height+y0)
	img = (clone.AsRGBA(img)).SubImage(rect)

	return img
}

// Resize takes an input image, width and height and returns the re-sized image
func (bp *BildProcessor) Resize(img image.Image, width, height int) image.Image {

	initW := img.Bounds().Dx()
	initH := img.Bounds().Dy()

	w, h := getResizeWidthAndHeight(width, height, initW, initH)
	if w != initW || h != initH {
		img = transform.Resize(img, w, h, transform.Linear)
	}

	return img
}

// Scale takes an input image, width and height and returns the re-sized
// image without maintaining the original aspect ratio
func (bp *BildProcessor) Scale(img image.Image, width, height int) image.Image {
	return transform.Resize(img, width, height, transform.Linear)
}

// GrayScale takes an input image and returns the grayscaled image
func (bp *BildProcessor) GrayScale(img image.Image) image.Image {
	// Rec. 601 Luma formula (https://en.wikipedia.org/wiki/Luma_%28video%29#Rec._601_luma_versus_Rec._709_luma_coefficients)
	return effect.GrayscaleWithWeights(img, 0.299, 0.587, 0.114)
}

// Blur takes an input image and blur radius and returns the Gausian blurred image
func (bp *BildProcessor) Blur(img image.Image, radius float64) image.Image {
	return blur.Gaussian(img, radius)
}

// Flip takes an input image and returns the image flipped. The direction of flip
// is determined by the specified mode - 'v' for a vertical flip, 'h' for a
// horizontal flip and 'vh'(or 'hv') for both.
func (bp *BildProcessor) Flip(img image.Image, mode string) image.Image {
	mode = strings.ToLower(mode)
	for _, op := range mode {
		switch op {
		case 'v':
			img = transform.FlipV(img)
		case 'h':
			img = transform.FlipH(img)
		}
	}
	return img
}

// Rotate takes an input image and returns a image rotated by the specified degrees.
// The rotation is applied clockwise, and fractional angles are also supported.
func (bp *BildProcessor) Rotate(img image.Image, angle float64) image.Image {
	return transform.Rotate(img, angle, nil)
}

// Decode takes a byte array and returns the decoded image, format, or the error
func (bp *BildProcessor) Decode(data []byte) (image.Image, string, error) {
	img, f, err := image.Decode(bytes.NewReader(data))
	return img, f, err
}

// Encode takes an image and the preferred format (extension) of the output
// Current supported format are "png", "jpg" and "jpeg"
func (bp *BildProcessor) Encode(img image.Image, fmt string) ([]byte, error) {
	enc := bp.encoders.GetEncoder(img, fmt)
	data, err := enc.Encode(img)
	return data, err
}

// FixOrientation takes an image and it's EXIF orientation
// To get the orientation of the image see GetOrientation (exif.go)
func (bp *BildProcessor) FixOrientation(img image.Image, orientation int) image.Image {
	switch orientation {
	case 2:
		return transform.FlipH(img)
	case 3:
		return transform.Rotate(img, 180, nil)
	case 4:
		img = transform.FlipH(img)
		return transform.Rotate(img, 180, nil)
	case 5:
		img = transform.FlipV(img)
		return transform.Rotate(img, 90, resizeBoundOption)
	case 6:
		return transform.Rotate(img, 90, resizeBoundOption)
	case 7:
		img = transform.FlipV(img)
		return transform.Rotate(img, 270, resizeBoundOption)
	case 8:
		return transform.Rotate(img, 270, resizeBoundOption)
	default:
		return img
	}
}

type overlayResult struct {
	overlayImg image.Image
	offset     image.Point
	index      int
	err        error
}

func (bp *BildProcessor) transformOverlay(i, w, h int, oa *processor.OverlayAttrs, c *chan overlayResult) {
	overlayImg, _, err := bp.Decode(oa.Img)
	fmt.Print(overlayImg)
	if err != nil {
		*c <- overlayResult{index: i, err: err}
	}
	if overlayImg == nil {
		*c <- overlayResult{index: i, err: fmt.Errorf("overlay byte cannot be decoded into image")}
	}

	ratio := float64(overlayImg.Bounds().Dy()) / float64(overlayImg.Bounds().Dx())
	dWidth := float64(w) * (oa.WidthPercentage / 100.0)

	// Resizing overlay image according to base image
	overlayImg = transform.Resize(overlayImg, int(dWidth), int(dWidth*ratio), transform.Linear)

	// Anchor point for overlaying
	x, y := getStartingPointForCrop(w, h, overlayImg.Bounds().Dx(), overlayImg.Bounds().Dy(), oa.Point)
	offset := image.Pt(x, y)
	*c <- overlayResult{
		overlayImg: overlayImg,
		offset:     offset,
		index:      i,
		err:        nil,
	}
}

// Watermark takes an input byte array, overlay byte array and opacity value
// and returns the watermarked image bytes or error
func (bp *BildProcessor) Watermark(base []byte, overlay []byte, opacity uint8) ([]byte, error) {
	baseImg, f, err := bp.Decode(base)
	if err != nil {
		return nil, err
	}
	if f != processor.ExtensionPNG {
		baseImg = clone.AsRGBA(baseImg)
	}

	oa := processor.OverlayAttrs{
		Img:              overlay,
		Point:            processor.PointCenter,
		WidthPercentage:  50.0,
		HeightPercentage: 50.0,
	}
	c := make(chan overlayResult)
	w := baseImg.Bounds().Dx()
	h := baseImg.Bounds().Dy()
	go bp.transformOverlay(0, w, h, &oa, &c)
	cr := <-c

	if cr.err != nil {
		return nil, cr.err
	}

	// Mask image (that is just a solid light gray image)
	mask := image.NewUniform(color.Alpha{A: opacity})

	// Performing overlay
	draw.DrawMask(baseImg.(draw.Image), cr.overlayImg.Bounds().Add(cr.offset), cr.overlayImg, image.ZP, mask, image.ZP, draw.Over)

	return bp.Encode(baseImg, f)
}

// Overlay takes a base image and array of overlay images and returns the final overlayed image bytes or error
func (bp *BildProcessor) Overlay(base []byte, overlays []*processor.OverlayAttrs) ([]byte, error) {
	if len(overlays) == 0 {
		return base, nil
	}

	baseImg, f, err := bp.Decode(base)
	if err != nil {
		return nil, err
	}
	if f != processor.ExtensionPNG {
		baseImg = clone.AsRGBA(baseImg)
	}

	c := make(chan overlayResult, len(overlays))
	w := baseImg.Bounds().Dx()
	h := baseImg.Bounds().Dy()
	for i, overlay := range overlays {
		go bp.transformOverlay(i, w, h, overlay, &c)
	}

	for i := 0; i < len(overlays); i++ {
		cr := <-c
		if cr.err == nil {
			// Performing overlay
			draw.DrawMask(baseImg.(draw.Image), cr.overlayImg.Bounds().Add(cr.offset), cr.overlayImg, image.ZP, nil, image.ZP, draw.Over)
		}
	}

	return bp.Encode(baseImg, f)
}

// WithEncoders is a builder function to set custom Encoders for BildProcessor
func WithEncoders(encoders *Encoders) ProcessorOption {
	return func(bp *BildProcessor) {
		bp.encoders = encoders
	}
}

// NewBildProcessor creates a new BildProcessor, if called without parameters encoders will be default
func NewBildProcessor(opts ...ProcessorOption) *BildProcessor {
	bp := &BildProcessor{encoders: NewEncoders()}
	for _, opt := range opts {
		opt(bp)
	}
	return bp
}
