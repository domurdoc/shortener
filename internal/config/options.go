package config

type Options struct {
	Addr            NetAddress
	BaseURL         URL
	LogLevel        LogLevel
	FileStoragePath FilePath
	DatabaseDSN     DataSourceName
}

func New(
	addr,
	baseURL,
	logLevel,
	storagePath,
	databaseDSN string,
) *Options {
	options := Options{}
	setOptionFromString(&options.BaseURL, baseURL)
	setOptionFromString(&options.Addr, addr)
	setOptionFromString(&options.LogLevel, logLevel)
	setOptionFromString(&options.FileStoragePath, storagePath)
	setOptionFromString(&options.DatabaseDSN, databaseDSN)
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
