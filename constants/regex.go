package constants

import "regexp"

var (
	S3Matcher         = regexp.MustCompile("^(?i)aws|s3$")
	CloudfrontMatcher = regexp.MustCompile("^(?i)cloudfront$")
	WebFolderMatcher  = regexp.MustCompile("^(?i)webfolder$")
)
