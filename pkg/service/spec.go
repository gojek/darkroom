package service

type Spec struct {
	// Scope defines a scope for the image manipulation job, it can be used for logging/mertrics collection purposes
	Scope string
	// ImageData holds the actual image contents to processed
	ImageData []byte
	// Params hold the key-value pairs for the processing job and tells the manipulator what to do with the image
	Params map[string]string
	// AcceptedFormats have the information whether darkroom can return the image using webp or not
	AcceptedFormats []string
}
