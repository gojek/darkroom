package config

type loggerConfig struct {
	level  string
	format string
}

type appInfo struct {
	name        string
	version     string
	description string
}

type bucketInfo struct {
	name       string
	accessKey  string
	secretKey  string
	pathPrefix string
}
