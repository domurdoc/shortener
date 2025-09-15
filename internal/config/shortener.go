package config

func LoadOptions() *Options {
	options := New(
		":8080",
		"http://localhost:8080",
		"info",
		"db.json",
		"user=yandex host=localhost port=5432 database=yandex sslmode=disable",
	)
	parseArgs(options)
	parseEnv(options)
	return options
}
