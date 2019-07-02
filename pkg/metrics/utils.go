package metrics

// GetImageSizeCluster takes in byte array and return the size cluster for tracking purpose
func GetImageSizeCluster(imageData []byte) string {
	switch sz := len(imageData); {
	case sz <= 128*1024:
		return "<=128KB"
	case sz <= 256*1024:
		return "<=256KB"
	case sz <= 512*1024:
		return "<=512KB"
	case sz <= 1024*1024:
		return "<=1MB"
	case sz <= 2048*1024:
		return "<=2MB"
	default:
		return ">2MB"
	}
}
