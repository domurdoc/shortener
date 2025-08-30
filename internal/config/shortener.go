package config

func LoadOptions() *Options {
	options := NewOptions(
		":8080",
		"http://localhost:8080",
	)
	parseArgs(options)
	parseEnv(options)
	return options
}
