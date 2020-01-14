package processor

// Point specifies which focus point in the image should be considered while cropping
type Point int

const (
	// PointTopLeft crops an image with focus point at top-left
	PointTopLeft Point = 1
	// PointTop crops an image with focus point at top
	PointTop Point = 2
	// PointTopRight crops an image with focus point at top-right
	PointTopRight Point = 3
	// PointLeft crops an image with focus point at left
	PointLeft Point = 4
	// PointCenter crops an image with focus point at center
	PointCenter Point = 5
	// PointRight crops an image with focus point at right
	PointRight Point = 6
	// PointBottomLeft crops an image with focus point at bottom-left
	PointBottomLeft Point = 7
	// PointBottom crops an image with focus point at bottom
	PointBottom Point = 8
	// PointBottomRight crops an image with focus point at bottom-right
	PointBottomRight Point = 9

	ExtensionWebP = "webp"
	ExtensionPNG  = "png"
	ExtensionJPG  = "jpg"
	ExtensionJPEG = "jpeg"
)
