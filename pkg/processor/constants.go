package processor

// CropPoint specifies which focus point in the image should be considered while cropping
type CropPoint int

const (
	// CropTopLeft crops an image with focus point at top-left
	CropTopLeft CropPoint = 1
	// CropTop crops an image with focus point at top
	CropTop CropPoint = 2
	// CropTopRight crops an image with focus point at top-right
	CropTopRight CropPoint = 3
	// CropLeft crops an image with focus point at left
	CropLeft CropPoint = 4
	// CropCenter crops an image with focus point at center
	CropCenter CropPoint = 5
	// CropRight crops an image with focus point at right
	CropRight CropPoint = 6
	// CropBottomLeft crops an image with focus point at bottom-left
	CropBottomLeft CropPoint = 7
	// CropBottom crops an image with focus point at bottom
	CropBottom CropPoint = 8
	// CropBottomRight crops an image with focus point at bottom-right
	CropBottomRight CropPoint = 9

	FormatJPG  = "jpg"
	FormatJPEG = "jpeg"
	FormatPNG  = "png"
	FormatWebP = "webp"
)
