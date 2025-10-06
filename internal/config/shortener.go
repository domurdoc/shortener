package config

func LoadOptions() *Options {
	options := New(
		":8080",
		"http://localhost:8080",
		"info",
		"",
		"",
		"ilovesber",
		"600s",
		"ilovesber",
		"600s",
		"5s",
	)
	parseArgs(options)
	parseEnv(options)
	return options
}
