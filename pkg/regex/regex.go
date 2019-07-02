package regex

import "regexp"

var (
	// S3Matcher regex matches against strings [aws|s3] in any case
	S3Matcher         = regexp.MustCompile("^(?i)aws|s3$")
	// CloudfrontMatcher regex matches against string cloudfront in any case
	CloudfrontMatcher = regexp.MustCompile("^(?i)cloudfront$")
	// WebFolderMatcher regex matches against string webfolder in any case
	WebFolderMatcher  = regexp.MustCompile("^(?i)webfolder$")
)
