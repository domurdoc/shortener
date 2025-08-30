package config

import "os"

func parseEnv(options *Options) {
	if value, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		if err := options.Addr.Set(value); err != nil {
			panic(err)
		}
	}
	if value, ok := os.LookupEnv("BASE_URL"); ok {
		if err := options.BaseURL.Set(value); err != nil {
			panic(err)
		}
	}
	if value, ok := os.LookupEnv("LOG_LEVEL"); ok {
		if err := options.LogLevel.Set(value); err != nil {
			panic(err)
		}
	}
}
