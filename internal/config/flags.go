package config

import "flag"

func parseArgs(options *Options) {
	flag.Var(&options.Addr, "a", "bind address")
	flag.Var(&options.BaseURL, "b", "base address")
	flag.Var(&options.LogLevel, "l", "logging level")
	flag.Var(&options.FileStoragePath, "f", "file storage path")
	flag.Var(&options.DatabaseDSN, "d", "database DSN")
	flag.Var(&options.DeleterMaxWorkers, "w", "deleter max workers")
	flag.Var(&options.DeleterMaxBatchSize, "s", "deleter max batch size")
	flag.Var(&options.DeleterCheckInterval, "c", "deleter check interval")
	flag.Parse()
}
