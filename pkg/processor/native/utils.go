package native

import (
	"image"

	"github.com/anthonynsimon/bild/parallel"
	"github.com/gojek/darkroom/pkg/config"
	"github.com/gojek/darkroom/pkg/processor"
)

func hasFastIsOpaque(im image.Image) bool {
	if _, ok := im.(*image.Gray); ok {
		return true
	}
	if _, ok := im.(*image.Gray16); ok {
		return true
	}
	if _, ok := im.(*image.CMYK); ok {
		return true
	}
	return false
}

func isOpaque(im image.Image) bool {
	// Check if image has fast Opaque checking method
	if hasFastIsOpaque(im) {
		oim, _ := im.(interface {
			Opaque() bool
		})
		return oim.Opaque()
	}
	// No fast Opaque() method, we need to loop through all pixels and check manually:
	rect := im.Bounds()
	isOpaque := true
	f := func(start, end int) {
		for y := rect.Min.Y + start; isOpaque && y < rect.Min.Y+end; y++ {
			for x := rect.Min.X; isOpaque && x < rect.Max.X; x++ {
				if _, _, _, a := im.At(x, y).RGBA(); a != 0xffff {
					isOpaque = false // Found a non-opaque pixel: image is non-opaque
				}
			}
		}
	}
	if config.ConcurrentOpacityCheckingEnabled() {
		parallel.Line(rect.Dy(), f)
	} else {
		f(rect.Min.Y, rect.Max.Y)
	}
	return isOpaque
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
func getStartingPointForCrop(w, h, rw, rh int, cropPoint processor.Point) (int, int) {
	x := (w - rw) / 2
	y := (h - rh) / 2

	switch cropPoint {
	case processor.PointTop:
		y = 0
	case processor.PointTopLeft:
		x = 0
		y = 0
	case processor.PointTopRight:
		x = w - rw
		y = 0
	case processor.PointLeft:
		x = 0
	case processor.PointRight:
		x = w - rw
	case processor.PointBottom:
		y = h - rh
	case processor.PointBottomLeft:
		x = 0
		y = h - rh
	case processor.PointBottomRight:
		x = w - rw
		y = h - rh
	}
	return x, y
}
