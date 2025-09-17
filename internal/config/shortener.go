package config

func LoadOptions() *Options {
	options := New(
		":8080",
		"http://localhost:8080",
		"info",
		"",
		"",
	)
	parseArgs(options)
	parseEnv(options)
	return options
}
