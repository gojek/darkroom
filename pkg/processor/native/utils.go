package native

import (
	"github.com/gojek/darkroom/pkg/processor"
	"image"
)

func isOpaque(im image.Image) bool {
	// Check if image has Opaque() method:
	if oim, ok := im.(interface {
		Opaque() bool
	}); ok {
		return oim.Opaque() // It does, call it and return its result!
	}
	// No Opaque() method, we need to loop through all pixels and check manually:
	rect := im.Bounds()
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			if _, _, _, a := im.At(x, y).RGBA(); a != 0xffff {
				return false // Found a non-opaque pixel: image is non-opaque
			}
		}
	}
	return true // All pixels are opaque, so is the image
}

// rw: required width, rh: required height, aw: actual width, ah: actual height
func getResizeWidthAndHeight(rw, rh, aw, ah int) (int, int) {
	if rh == 0 {
		h := (rw * ah) / aw
		return rw, h
	} else if rw == 0 {
		w := (rh * aw) / ah
		return w, rh
	} else {
		h := (rw * ah) / aw
		if h <= rh {
			return rw, h
		}
		w := (rh * aw) / ah
		return w, rh
	}
}

// rw: required width, rh: required height, aw: actual width, ah: actual height
func getResizeWidthAndHeightForCrop(rw, rh, aw, ah int) (int, int) {
	h := (rw * ah) / aw
	if h >= rh {
		return rw, h
	}
	w := (rh * aw) / ah
	return w, rh
}

// w: scaled width, h: scaled height, rw: required width, rh: required height
func getStartingPointForCrop(w, h, rw, rh int, cropPoint processor.CropPoint) (int, int) {
	x := (w - rw) / 2
	y := (h - rh) / 2

	switch cropPoint {
	case processor.CropTop:
		y = 0
		break
	case processor.CropTopLeft:
		x = 0
		y = 0
		break
	case processor.CropTopRight:
		x = w - rw
		y = 0
		break
	case processor.CropLeft:
		x = 0
		break
	case processor.CropRight:
		x = w - rw
		break
	case processor.CropBottom:
		y = h - rh
		break
	case processor.CropBottomLeft:
		x = 0
		y = h - rh
		break
	case processor.CropBottomRight:
		x = w - rw
		y = h - rh
		break
	default:
		break
	}
	return x, y
}
