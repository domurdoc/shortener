package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/caarlos0/env/v11"

	"github.com/domurdoc/shortener/internal/utils"
)

type BaseConfig struct {
	ConfigFile    string `env:"CONFIG" json:"-"`
	LoggerLevel   string `env:"LOG_LEVEL" json:"log_level"`
	TrustedSubnet string `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
}

type ServerConfig struct {
	ServerAddress      string        `env:"SERVER_ADDRESS" json:"server_address"`
	ServerCloseTimeout time.Duration `env:"SERVER_CLOSE_TIMEOUT" json:"server_close_timeout"`
	ServerEnableHTTPS  bool          `env:"ENABLE_HTTPS" json:"enable_https"`
	ServerCertFile     string        `env:"CERT_FILE" json:"cert_file"`
	ServerKeyFile      string        `env:"KEY_FILE" json:"key_file"`
}

type AuthConfig struct {
	AuthJWTSecret    string        `env:"JWT_SECRET" json:"jwt_secret"`
	AuthJWTDuration  time.Duration `env:"JWT_DURATION" json:"jwt_duration"`
	AuthCookieName   string        `env:"COOKIE_NAME" json:"cookie_name"`
	AuthCookieMaxAge time.Duration `env:"COOKIE_MAX_AGE" json:"cookie_max_age"`
}

type AuditConfig struct {
	AuditFilePath          string        `env:"AUDIT_FILE" json:"audit_file"`
	AuditFilePoolSize      int           `env:"AUDIT_FILE_POOL_SIZE" json:"audit_file_pool_size"`
	AuditFileMaxBatchSize  int           `env:"AUDIT_FILE_MAX_BATCH_SIZE" json:"audit_file_max_batch_size"`
	AuditFileBatchInterval time.Duration `env:"AUDIT_FILE_BATCH_INTERVAL" json:"audit_file_batch_interval"`
	AuditRemoteURL         string        `env:"AUDIT_URL" json:"audit_url"`
	AuditRemotePoolSize    int           `env:"AUDIT_REMOTE_POOL_SIZE" json:"audit_remote_pool_size"`
}

type ServiceConfig struct {
	ServiceBaseURL                string        `env:"BASE_URL" json:"base_url"`
	ServiceDeleterMaxWorkers      int           `env:"DELETER_MAX_WORKERS" json:"deleter_max_workers"`
	ServiceDeleterMaxBatchSize    int           `env:"DELETER_MAX_BATCH_SIZE" json:"deleter_max_batch_size"`
	ServiceDeleterCheckInterval   time.Duration `env:"DELETER_CHECK_INTERVAL" json:"deleter_check_interval"`
	ServiceGeneratorRandomLength  int           `env:"RANDOM_CODE_LENGTH" json:"random_code_length"`
	ServiceGeneratorRandomCharSet string        `env:"RANDOM_CODE_CHARSET" json:"random_code_charset"`
	ServiceGeneratorConstantValue string        `env:"CONSTANT_CODE_VALUE" json:"constant_code_value"`
}

type RepositoriesConfig struct {
	RepositoryDSN      string `env:"DATABASE_DSN" json:"database_dsn"`
	RepositoryFilePath string `env:"FILE_STORAGE_PATH" json:"file_storage_path"`
}

type ProfilerConfig struct {
	ProfilerAddress      string        `env:"PROFILER_ADDRESS" json:"profiler_address"`
	ProfilerCloseTimeout time.Duration `env:"PROFILER_CLOSE_TIMEOUT" json:"profiler_close_timeout"`
}

type LoggerConfig struct {
	LoggerLevel string `env:"LOG_LEVEL" json:"log_level"`
}

type GRPCConfig struct {
	GRPCPort int `env:"GRPC_PORT" json:"grpc_port"`
}

// Config holds the complete configuration for the shortener application.
// It is populated from environment variables, command-line flags, and default values.
// The struct uses nested groups to organize settings by concern (server, auth, audit, etc.).
type Config struct {
	BaseConfig
	ServerConfig
	AuthConfig
	AuditConfig
	ServiceConfig
	RepositoriesConfig
	ProfilerConfig
	GRPCConfig
}

func Default() *Config {
	cfg := &Config{}
	cfg.ServerAddress = "localhost:8080"
	cfg.ServerCloseTimeout = 10 * time.Second
	cfg.ServerEnableHTTPS = false
	cfg.AuthJWTDuration = 600 * time.Second
	cfg.AuthJWTSecret = utils.MustGenerateRandomString(utils.ALPHA, 32)
	cfg.AuthCookieName = "ilovesber"
	cfg.AuthCookieMaxAge = 600 * time.Second
	cfg.AuditFilePoolSize = 1
	cfg.AuditFileMaxBatchSize = 10
	cfg.AuditFileBatchInterval = 1 * time.Second
	cfg.AuditRemotePoolSize = 1
	cfg.LoggerLevel = "debug"
	cfg.ServiceBaseURL = "http://localhost:8080"
	cfg.ServiceDeleterMaxWorkers = 2
	cfg.ServiceDeleterMaxBatchSize = 10
	cfg.ServiceDeleterCheckInterval = 1 * time.Second
	cfg.ServiceGeneratorRandomLength = 6
	cfg.ServiceGeneratorRandomCharSet = utils.ALPHA
	cfg.ProfilerCloseTimeout = 10 * time.Second
	cfg.TrustedSubnet = ""
	cfg.GRPCPort = 0
	return cfg
}

func ParseEnv(cfg *Config) (*Config, error) {
	return cfg, env.Parse(cfg)
}

func ParseArgs(cfg *Config) (*Config, error) {
	flag.StringVar(&cfg.ConfigFile, "c", cfg.ConfigFile, "config file") // for usage purposes only
	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "bind address")
	flag.BoolVar(&cfg.ServerEnableHTTPS, "s", cfg.ServerEnableHTTPS, "enable https")
	flag.StringVar(&cfg.ServerCertFile, "tls-cert", cfg.ServerCertFile, "cert file")
	flag.StringVar(&cfg.ServerKeyFile, "tls-key", cfg.ServerKeyFile, "key file")
	flag.StringVar(&cfg.ServiceBaseURL, "b", cfg.ServiceBaseURL, "base address")
	flag.StringVar(&cfg.LoggerLevel, "l", cfg.LoggerLevel, "logging level")
	flag.StringVar(&cfg.RepositoryFilePath, "f", cfg.RepositoryFilePath, "file storage path")
	flag.StringVar(&cfg.RepositoryDSN, "d", cfg.RepositoryDSN, "database DSN")
	flag.IntVar(&cfg.ServiceDeleterMaxWorkers, "w", cfg.ServiceDeleterMaxWorkers, "deleter max workers")
	flag.IntVar(&cfg.ServiceDeleterMaxBatchSize, "ds", cfg.ServiceDeleterMaxBatchSize, "deleter max batch size")
	flag.DurationVar(&cfg.ServiceDeleterCheckInterval, "i", cfg.ServiceDeleterCheckInterval, "deleter check interval")
	flag.StringVar(&cfg.AuditFilePath, "audit-file", cfg.AuditFilePath, "audit file")
	flag.StringVar(&cfg.AuditRemoteURL, "audit-url", cfg.AuditRemoteURL, "audit url")
	flag.StringVar(&cfg.ProfilerAddress, "p", cfg.ProfilerAddress, "pprof address")
	flag.StringVar(&cfg.TrustedSubnet, "t", cfg.TrustedSubnet, "trusted subnet")
	flag.IntVar(&cfg.GRPCPort, "g", cfg.GRPCPort, "grpc port")
	flag.Parse()
	return cfg, nil
}

func ParseFile(cfg *Config, path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	dec := json.NewDecoder(file)
	if err := dec.Decode(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file (JSON): %w", err)
	}
	return cfg, nil
}

func LoadConfig() (*Config, error) {
	var err error

	cfg := Default()
	configPath := getConfigPath()
	if configPath != "" {
		cfg, err = ParseFile(cfg, configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to parse json file: %w", err)
		}
	}
	cfg, err = ParseEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse env: %w", err)
	}
	cfg, err = ParseArgs(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse args: %w", err)
	}
	return cfg, nil
}

func getConfigPath() string {
	var configPath string

	i := slices.Index(os.Args, "-c")
	if i != -1 && i+1 < len(os.Args) {
		configPath = os.Args[i+1]
	}
	if configPath == "" {
		configPath = os.Getenv("CONFIG")
	}
	return configPath
}
