package gpu

import (
	"github.com/gojek/darkroom/pkg/processor"
	"github.com/gojek/darkroom/pkg/processor/utils"
	"gocv.io/x/gocv"
	"gocv.io/x/gocv/cuda"
	"image"
)

type OpenCVProcessor struct {
}

func (p *OpenCVProcessor) Crop(image image.Image, width, height int, point processor.CropPoint) image.Image {
	panic("implement me")
}

func (p *OpenCVProcessor) Resize(image image.Image, width, height int) image.Image {
	panic("implement me")
}

func (p *OpenCVProcessor) GpuResize(source []byte, width, height int) []byte {
	src, err := gocv.IMDecode(source, gocv.IMReadColor)
	if err != nil {
		return source
	}
	defer src.Close()

	initW := src.Cols()
	initH := src.Rows()
	w, h := utils.GetResizeWidthAndHeight(width, height, initW, initH)
	if w != initW || h != initH {
		gMatSrc, gMatDst := cuda.NewGpuMat(), cuda.NewGpuMat()
		defer func() {
			_ = gMatSrc.Close()
			_ = gMatDst.Close()
		}()

		gMatSrc.Upload(src)

		dst := gocv.NewMat()
		defer dst.Close()

		cuda.Resize(gMatSrc, &gMatDst, image.Point{X: w, Y: h}, 0, 0, gocv.InterpolationCubic)

		gMatDst.Download(&dst)

		data, err := gocv.IMEncode(gocv.JPEGFileExt, dst)
		if err != nil {
			return source
		}
		return data
	}
	return source
}

func (p *OpenCVProcessor) GrayScale(image image.Image) image.Image {
	panic("implement me")
}

func (p *OpenCVProcessor) Blur(image image.Image, radius float64) image.Image {
	panic("implement me")
}

func (p *OpenCVProcessor) Watermark(base []byte, overlay []byte, opacity uint8) ([]byte, error) {
	panic("implement me")
}

func (p *OpenCVProcessor) Flip(image image.Image, mode string) image.Image {
	panic("implement me")
}

func (p *OpenCVProcessor) Rotate(image image.Image, angle float64) image.Image {
	panic("implement me")
}

func (p *OpenCVProcessor) Decode(data []byte) (img image.Image, format string, err error) {
	panic("implement me")
}

func (p *OpenCVProcessor) Encode(img image.Image, format string) ([]byte, error) {
	panic("implement me")
}

func (p *OpenCVProcessor) FixOrientation(img image.Image, orientation int) image.Image {
	panic("implement me")
}

func NewOpenCVProcessor() *OpenCVProcessor {
	return &OpenCVProcessor{}
}
