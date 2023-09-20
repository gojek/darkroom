package service

type ProcessSpec interface {
	// IsWebPSupported() will tell if WebP is supported based on the accepted formats
	IsWebPSupported() bool
}

type processSpec struct {
	// Scope defines a scope for the image manipulation job, it can be used for logging/mertrics collection purposes
	Scope string
	// ImageData holds the actual image contents to processed
	ImageData []byte
	// Params hold the key-value pairs for the processing job and tells the manipulator what to do with the image
	Params map[string]string
	// Reformat image to target format
	TargetFormat string
	// Formats have the information of accepted formats, whether darkroom can return the image using webp or not
	formats []string
}

const (
	extJPG  = "jpg"
	extPNG  = "png"
	extWebP = "webp"
	extJPEG = "jpeg"
)

func (ps *processSpec) IsWebPSupported() bool {
	for _, f := range ps.formats {
		if f == "image/webp" {
			return true
		}
	}
	return false
}

type SpecBuilder interface {
	WithScope(scope string) SpecBuilder
	WithImageData(img []byte) SpecBuilder
	WithParams(params map[string]string) SpecBuilder
	WithFormats(formats []string) SpecBuilder
	WithTargetFormat(ext string) SpecBuilder
	Build() processSpec
}

type specBuilder struct {
	scope     string
	imageData []byte
	params    map[string]string
	formats   []string
	extension string
}

func (sb *specBuilder) WithScope(scope string) SpecBuilder {
	sb.scope = scope
	return sb
}

func (sb *specBuilder) WithImageData(img []byte) SpecBuilder {
	sb.imageData = img
	return sb
}

func (sb *specBuilder) WithParams(params map[string]string) SpecBuilder {
	sb.params = params
	return sb
}

func (sb *specBuilder) WithFormats(formats []string) SpecBuilder {
	sb.formats = formats
	return sb
}

func (sb *specBuilder) WithTargetFormat(ext string) SpecBuilder {
	switch ext {
	case extJPG, extJPEG, extPNG, extWebP:
		sb.extension = ext
	}
	return sb
}

func (sb *specBuilder) Build() processSpec {
	return processSpec{
		Scope:        sb.scope,
		ImageData:    sb.imageData,
		Params:       sb.params,
		formats:      sb.formats,
		TargetFormat: sb.extension,
	}
}

func NewSpecBuilder() SpecBuilder {
	return &specBuilder{}
}
