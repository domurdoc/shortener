package config

import (
	"fmt"
)

type Options struct {
	Addr            NetAddress
	BaseURL         URL
	LogLevel        LogLevel
	FileStoragePath FilePath
}

func (o Options) String() string {
	return fmt.Sprintf(
		"addr = %q; baseURL = %q; logLevel = %q; fileStoragePath = %q",
		o.Addr,
		o.BaseURL,
		o.LogLevel,
		o.FileStoragePath,
	)
}

func New(defaultAddr, defaultBaseURL, defaultLogLevel, defaultStoragePath string) *Options {
	options := Options{}
	setOptionFromString(&options.BaseURL, defaultBaseURL)
	setOptionFromString(&options.Addr, defaultAddr)
	setOptionFromString(&options.LogLevel, defaultLogLevel)
	setOptionFromString(&options.FileStoragePath, defaultStoragePath)
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
