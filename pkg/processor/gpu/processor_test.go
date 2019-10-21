package gpu

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"image/jpeg"
	"io/ioutil"
	"testing"
)

type OpenCVProcessorSuite struct {
	suite.Suite
	srcData       []byte
	processor     *OpenCVProcessor
}

func (s *OpenCVProcessorSuite) SetupSuite() {
	s.processor = NewOpenCVProcessor()
	s.srcData, _ = ioutil.ReadFile("../_testdata/test.jpg")
}

func TestOpenCVProcessor(t *testing.T) {
	suite.Run(t, new(OpenCVProcessorSuite))
}

func (s *OpenCVProcessorSuite) TestOpenCVProcessor_GpuResize() {
	out := s.processor.GpuResize(s.srcData, 600, 500)

	img, _ := jpeg.Decode(bytes.NewReader(out))
	assert.NotNil(s.T(), out)
	assert.Equal(s.T(), 600, img.Bounds().Dx())
	assert.Equal(s.T(), 450, img.Bounds().Dy())
}

func BenchmarkOpenCVProcessor_GpuResize(b *testing.B) {
	p := NewOpenCVProcessor()
	img, _ := ioutil.ReadFile("../_testdata/test.jpg")
	for n := 0; n < b.N; n++ {
		_ = p.GpuResize(img, 3000, 4000)
	}
}
