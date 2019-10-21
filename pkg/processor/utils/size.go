package utils

import "github.com/gojek/darkroom/pkg/processor"

// rw: required width, rh: required height, aw: actual width, ah: actual height
func GetResizeWidthAndHeight(rw, rh, aw, ah int) (int, int) {
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
func GetResizeWidthAndHeightForCrop(rw, rh, aw, ah int) (int, int) {
	h := (rw * ah) / aw
	if h >= rh {
		return rw, h
	}
	w := (rh * aw) / ah
	return w, rh
}

// w: scaled width, h: scaled height, rw: required width, rh: required height
func GetStartingPointForCrop(w, h, rw, rh int, cropPoint processor.CropPoint) (int, int) {
	x := (w - rw) / 2
	y := (h - rh) / 2

	switch cropPoint {
	case processor.CropTop:
		y = 0
	case processor.CropTopLeft:
		x = 0
		y = 0
	case processor.CropTopRight:
		x = w - rw
		y = 0
	case processor.CropLeft:
		x = 0
	case processor.CropRight:
		x = w - rw
	case processor.CropBottom:
		y = h - rh
	case processor.CropBottomLeft:
		x = 0
		y = h - rh
	case processor.CropBottomRight:
		x = w - rw
		y = h - rh
	}
	return x, y
}
