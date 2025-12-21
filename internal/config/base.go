package config

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v11"

	"github.com/domurdoc/shortener/internal/utils"
)

// Config holds the complete configuration for the shortener application.
// It is populated from environment variables, command-line flags, and default values.
// The struct uses nested groups to organize settings by concern (server, auth, audit, etc.).
type Config struct {
	Server struct {
		Address      string        `env:"SERVER_ADDRESS"`
		CloseTimeout time.Duration `env:"SERVER_CLOSE_TIMEOUT"`
	}
	Auth struct {
		Strategy struct {
			JWTSecret   string        `env:"JWT_SECRET"`
			JWTDuration time.Duration `env:"JWT_DURATION"`
		}
		Transport struct {
			CookieName   string        `env:"COOKIE_NAME"`
			CookieMaxAge time.Duration `env:"COOKIE_MAX_AGE"`
		}
	}
	Audit struct {
		File struct {
			Path          string        `env:"AUDIT_FILE"`
			PoolSize      int           `env:"AUDIT_FILE_POOL_SIZE"`
			MaxBatchSize  int           `env:"AUDIT_FILE_MAX_BATCH_SIZE"`
			BatchInterval time.Duration `env:"AUDIT_FILE_BATCH_INTERVAL"`
		}
		Remote struct {
			URL      string `env:"AUDIT_URL"`
			PoolSize int    `env:"AUDIT_REMOTE_POOL_SIZE"`
		}
	}
	Logger struct {
		Level string `env:"LOG_LEVEL"`
	}
	Service struct {
		BaseURL              string        `env:"BASE_URL"`
		DeleterMaxWorkers    int           `env:"DELETER_MAX_WORKERS"`
		DeleterMaxBatchSize  int           `env:"DELETER_MAX_BATCH_SIZE"`
		DeleterCheckInterval time.Duration `env:"DELETER_CHECK_INTERVAL"`
	}
	Repositories struct {
		DB struct {
			DSN string `env:"DATABASE_DSN"`
		}
		File struct {
			Path string `env:"FILE_STORAGE_PATH"`
		}
	}
	Generator struct {
		Random struct {
			Length  int    `env:"RANDOM_CODE_LENGTH"`
			CharSet string `env:"RANDOM_CODE_CHARSET"`
		}
		Constant struct {
			Value string `env:"CONSTANT_CODE_VALUE"`
		}
	}
	Profiler struct {
		Address      string        `env:"PROFILER_ADDRESS"`
		CloseTimeout time.Duration `env:"PROFILER_CLOSE_TIMEOUT"`
	}
}

func New() *Config {
	cfg := &Config{}
	cfg.Server.Address = "localhost:8080"
	cfg.Server.CloseTimeout = 10 * time.Second
	cfg.Auth.Strategy.JWTDuration = 600 * time.Second
	cfg.Auth.Strategy.JWTSecret = utils.MustGenerateRandomString(utils.ALPHA, 32)
	cfg.Auth.Transport.CookieName = "ilovesber"
	cfg.Auth.Transport.CookieMaxAge = 600 * time.Second
	cfg.Audit.File.PoolSize = 1
	cfg.Audit.File.MaxBatchSize = 10
	cfg.Audit.File.BatchInterval = 1 * time.Second
	cfg.Audit.Remote.PoolSize = 1
	cfg.Logger.Level = "debug"
	cfg.Service.BaseURL = "http://localhost:8080"
	cfg.Service.DeleterMaxWorkers = 2
	cfg.Service.DeleterMaxBatchSize = 10
	cfg.Service.DeleterCheckInterval = 1 * time.Second
	cfg.Generator.Random.Length = 6
	cfg.Generator.Random.CharSet = utils.ALPHA
	cfg.Profiler.CloseTimeout = 10 * time.Second
	return cfg
}

func ParseEnv(cfg *Config) error {
	return env.Parse(cfg)
}

func ParseArgs(cfg *Config) {
	flag.StringVar(&cfg.Server.Address, "a", cfg.Server.Address, "bind address")
	flag.StringVar(&cfg.Service.BaseURL, "b", cfg.Service.BaseURL, "base address")
	flag.StringVar(&cfg.Logger.Level, "l", cfg.Logger.Level, "logging level")
	flag.StringVar(&cfg.Repositories.File.Path, "f", cfg.Repositories.File.Path, "file storage path")
	flag.StringVar(&cfg.Repositories.DB.DSN, "d", cfg.Repositories.DB.DSN, "database DSN")
	flag.IntVar(&cfg.Service.DeleterMaxWorkers, "w", cfg.Service.DeleterMaxWorkers, "deleter max workers")
	flag.IntVar(&cfg.Service.DeleterMaxBatchSize, "s", cfg.Service.DeleterMaxBatchSize, "deleter max batch size")
	flag.DurationVar(&cfg.Service.DeleterCheckInterval, "c", cfg.Service.DeleterCheckInterval, "deleter check interval")
	flag.StringVar(&cfg.Audit.File.Path, "audit-file", cfg.Audit.File.Path, "audit file")
	flag.StringVar(&cfg.Audit.Remote.URL, "audit-url", cfg.Audit.Remote.URL, "audit url")
	flag.StringVar(&cfg.Profiler.Address, "p", cfg.Profiler.Address, "pprof address")
	flag.Parse()
}
