package config

import "os"

func parseEnv(options *Options) {
	serverAddress, ok := os.LookupEnv("SERVER_ADDRESS")
	if ok {
		if err := options.Addr.Set(serverAddress); err != nil {
			panic(err)
		}
	}
	baseURL, ok := os.LookupEnv("BASE_URL")
	if ok {
		if err := options.BaseURL.Set(baseURL); err != nil {
			panic(err)
		}
	}
}
