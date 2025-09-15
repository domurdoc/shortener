package config

import "fmt"

type Options struct {
	Addr            NetAddress
	BaseURL         URL
	LogLevel        LogLevel
	FileStoragePath FilePath
	DatabaseDSN     DataSourceName
}

func (o Options) String() string {
	return fmt.Sprintf(
		"addr = %q; baseURL = %q; logLevel = %q; fileStoragePath = %q; databaseDSB = %q",
		o.Addr,
		o.BaseURL,
		o.LogLevel,
		o.FileStoragePath,
		o.DatabaseDSN,
	)
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
