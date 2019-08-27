package service

type Spec interface {
	// IsWebPSupported() will tell if WebP is supported based on the accepted formats
	IsWebPSupported() bool
}

type spec struct {
	// Scope defines a scope for the image manipulation job, it can be used for logging/mertrics collection purposes
	Scope string
	// ImageData holds the actual image contents to processed
	ImageData []byte
	// Params hold the key-value pairs for the processing job and tells the manipulator what to do with the image
	Params map[string]string
	// Formats have the information of accepted formats, whether darkroom can return the image using webp or not
	formats []string
}

func (s *spec) IsWebPSupported() bool {
	return true
}

type SpecBuilder interface {
	WithScope(scope string) SpecBuilder
	WithImageData(img []byte) SpecBuilder
	WithParams(params map[string]string) SpecBuilder
	WithFormats(formats []string) SpecBuilder
	Build() spec
}

type specBuilder struct {
	scope     string
	imageData []byte
	params    map[string]string
	formats   []string
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

func (sb *specBuilder) Build() spec {
	return spec{
		Scope:     sb.scope,
		ImageData: sb.imageData,
		Params:    sb.params,
		formats:   sb.formats,
	}
}

func NewSpecBuilder() SpecBuilder {
	return &specBuilder{}
}
