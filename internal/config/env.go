package config

import "os"

func parseEnv(options *Options) {
	setOptionFromEnv(&options.Addr, "SERVER_ADDRESS")
	setOptionFromEnv(&options.BaseURL, "BASE_URL")
	setOptionFromEnv(&options.LogLevel, "LOG_LEVEL")
	setOptionFromEnv(&options.FileStoragePath, "FILE_STORAGE_PATH")
	setOptionFromEnv(&options.JWTSecret, "JWT_SECRET")
	setOptionFromEnv(&options.JWTDuration, "JWT_DURATION")
	setOptionFromEnv(&options.CookieMaxAge, "COOKIE_MAX_AGE")
}

func setOptionFromEnv(s option, envName string) {
	if value, ok := os.LookupEnv(envName); ok {
		setOptionFromString(s, value)
	}
}
