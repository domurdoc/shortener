package config

import "os"

func parseEnv(options *Options) {
	setOptionFromEnv(&options.Addr, "SERVER_ADDRESS")
	setOptionFromEnv(&options.BaseURL, "BASE_URL")
	setOptionFromEnv(&options.LogLevel, "LOG_LEVEL")
	setOptionFromEnv(&options.FileStoragePath, "FILE_STORAGE_PATH")
}

func setOptionFromEnv(s option, envName string) {
	if value, ok := os.LookupEnv(envName); ok {
		setOptionFromString(s, value)
	}
}
