package config

import "flag"

func parseArgs(options *Options) {
	flag.Var(&options.Addr, "a", "bind address")
	flag.Var(&options.BaseURL, "b", "base address")
	flag.Var(&options.LogLevel, "l", "logging level")
	flag.Parse()
}
