package config

type Options struct {
	Addr                  NetAddress
	BaseURL               URL
	LogLevel              LogLevel
	FileStoragePath       String
	DatabaseDSN           String
	JWTSecret             String
	JWTDuration           Duration
	CookieName            String
	CookieMaxAge          Duration
	SaveDeletionsInterval Duration
}

func New(
	addr,
	baseURL,
	logLevel,
	storagePath,
	databaseDSN,
	jwtSecret,
	jwtDuration,
	cookieName,
	cookieMaxAge,
	saveDeletionsInterval string,
) *Options {
	options := Options{}
	setOptionFromString(&options.BaseURL, baseURL)
	setOptionFromString(&options.Addr, addr)
	setOptionFromString(&options.LogLevel, logLevel)
	setOptionFromString(&options.FileStoragePath, storagePath)
	setOptionFromString(&options.DatabaseDSN, databaseDSN)
	setOptionFromString(&options.JWTSecret, jwtSecret)
	setOptionFromString(&options.JWTDuration, jwtDuration)
	setOptionFromString(&options.CookieMaxAge, cookieMaxAge)
	setOptionFromString(&options.CookieName, cookieName)
	setOptionFromString(&options.SaveDeletionsInterval, saveDeletionsInterval)
	return &options
}

type option interface {
	Set(string) error
}

func setOptionFromString(o option, value string) {
	if err := o.Set(value); err != nil {
		panic(err)
	}
}
