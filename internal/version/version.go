package version

type BuildInfo struct {
	Version   string
	GitTag    string
	GitCommit string
	BuildDate string
}

var (
	Build BuildInfo
)
