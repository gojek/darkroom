package service

import (
	"context"
	"***REMOVED***/darkroom/processor"
	"strconv"
)

const (
	width        = "w"
	height       = "h"
	fit          = "fit"
	crop         = "crop"
	mono         = "mono"
	blackHexCode = "000000"
)

type Manipulator interface {
	Process(ctx context.Context, data []byte, params map[string]string) ([]byte, error)
}

type manipulator struct {
	processor processor.Processor
}

func (m *manipulator) Process(ctx context.Context, data []byte, params map[string]string) ([]byte, error) {
	var err error
	if params[fit] == crop {
		data, err = m.processor.Crop(data, CleanInt(params[width]), CleanInt(params[height]), GetCropPoint(params[crop]))
	} else if len(params[fit]) == 0 && (CleanInt(params[width]) != 0 || CleanInt(params[height]) != 0) {
		data, err = m.processor.Resize(data, CleanInt(params[width]), CleanInt(params[height]))
	}
	if params[mono] == blackHexCode {
		data, err = m.processor.GrayScale(data)
	}
	return data, err
}

func CleanInt(input string) int {
	val, _ := strconv.Atoi(input)
	if val <= 0 {
		return 0
	}
	return val % 10000 // Never return value greater than 9999
}

func GetCropPoint(input string) processor.CropPoint {
	switch input {
	case "top":
		return processor.CropTop
	case "top,left":
		return processor.CropTopLeft
	case "top,right":
		return processor.CropTopRight
	case "left":
		return processor.CropLeft
	case "right":
		return processor.CropRight
	case "bottom":
		return processor.CropBottom
	case "bottom,left":
		return processor.CropBottomLeft
	case "bottom,right":
		return processor.CropBottomRight
	default:
		return processor.CropCenter
	}
}

func NewManipulator(processor processor.Processor) *manipulator {
	return &manipulator{processor: processor}
}
