package config

import "github.com/domurdoc/shortener/internal/utils"

func LoadOptions() *Options {
	options := New(
		":8080",
		"http://localhost:8080",
		"info",
		"",
		"",
		utils.GenerateRandomString(utils.ALPHA, 32),
		"600s",
		"ilovesber",
		"600s",
		"10",
		"10",
		"5s",
		"",
		"1",
		"0",
		"2s",
		"",
		"1",
	)
	parseArgs(options)
	parseEnv(options)
	return options
}
