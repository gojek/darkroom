package processor

type CropPoint int

const (
	CropTopLeft     CropPoint = 1
	CropTop         CropPoint = 2
	CropTopRight    CropPoint = 3
	CropLeft        CropPoint = 4
	CropCenter      CropPoint = 5
	CropRight       CropPoint = 6
	CropBottomLeft  CropPoint = 7
	CropBottom      CropPoint = 8
	CropBottomRight CropPoint = 9
)
