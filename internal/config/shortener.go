package config

func LoadOptions() *Options {
	options := NewOptions(
		":8080",
		"http://localhost:8080",
		"info",
		"db.json",
	)
	parseArgs(options)
	parseEnv(options)
	return options
}
