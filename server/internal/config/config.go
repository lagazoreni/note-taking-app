package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultAddr           = "127.0.0.1:8080"
	DefaultDataDir        = "./.local/data"
	DefaultDBFile         = "notes.db"
	DefaultLogLevel       = "info"
	DefaultImportMaxBytes = int64(256 * 1024 * 1024)
	DefaultImportTTL      = 24 * time.Hour
	DefaultApplicationVer = "0.1.0"
)

type Config struct {
	Addr           string
	DataDir        string
	DBFile         string
	LogLevel       string
	ImportMaxBytes int64
	ImportTTL      time.Duration
	Version        string
}

func Load() (Config, error) { return FromEnvironment() }

func FromEnvironment() (Config, error) {
	env := make(map[string]string)
	for _, item := range os.Environ() {
		parts := strings.SplitN(item, "=", 2)
		env[parts[0]] = parts[1]
	}
	return FromEnv(env)
}

// FromEnv parses a map to make configuration deterministic and easy to test.
func Parse(env map[string]string) (Config, error) { return FromEnv(env) }

func FromEnv(env map[string]string) (Config, error) {
	cfg := Config{
		Addr: DefaultAddr, DataDir: DefaultDataDir, DBFile: DefaultDBFile,
		LogLevel: DefaultLogLevel, ImportMaxBytes: DefaultImportMaxBytes,
		ImportTTL: DefaultImportTTL, Version: DefaultApplicationVer,
	}
	if value, ok := env["NOTED_ADDR"]; ok {
		cfg.Addr = strings.TrimSpace(value)
	}
	if value, ok := env["NOTED_DATA_DIR"]; ok {
		cfg.DataDir = strings.TrimSpace(value)
	}
	if value, ok := env["NOTED_DB_FILE"]; ok {
		cfg.DBFile = strings.TrimSpace(value)
	}
	if value, ok := env["NOTED_LOG_LEVEL"]; ok {
		cfg.LogLevel = strings.ToLower(strings.TrimSpace(value))
	}
	if value, ok := env["NOTED_VERSION"]; ok && strings.TrimSpace(value) != "" {
		cfg.Version = strings.TrimSpace(value)
	}
	if value, ok := env["NOTED_IMPORT_MAX_BYTES"]; ok {
		parsed, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
		if err != nil {
			return Config{}, fmt.Errorf("NOTED_IMPORT_MAX_BYTES must be an integer: %w", err)
		}
		cfg.ImportMaxBytes = parsed
	}
	if value, ok := env["NOTED_IMPORT_TTL"]; ok {
		parsed, err := time.ParseDuration(strings.TrimSpace(value))
		if err != nil {
			return Config{}, fmt.Errorf("NOTED_IMPORT_TTL must be a duration: %w", err)
		}
		cfg.ImportTTL = parsed
	}
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c Config) Validate() error {
	if strings.TrimSpace(c.Addr) == "" {
		return fmt.Errorf("NOTED_ADDR must not be empty")
	}
	if strings.TrimSpace(c.DataDir) == "" {
		return fmt.Errorf("NOTED_DATA_DIR must not be empty")
	}
	if strings.TrimSpace(c.DBFile) == "" || filepath.Base(c.DBFile) != c.DBFile || c.DBFile == "." || c.DBFile == ".." {
		return fmt.Errorf("NOTED_DB_FILE must be a file name within NOTED_DATA_DIR")
	}
	if c.ImportMaxBytes <= 0 {
		return fmt.Errorf("NOTED_IMPORT_MAX_BYTES must be greater than zero")
	}
	if c.ImportTTL <= 0 {
		return fmt.Errorf("NOTED_IMPORT_TTL must be greater than zero")
	}
	switch c.LogLevel {
	case "debug", "info", "warn", "error":
	default:
		return fmt.Errorf("NOTED_LOG_LEVEL must be one of debug, info, warn, error")
	}
	return nil
}

func (c Config) DatabasePath() string { return filepath.Join(c.DataDir, c.DBFile) }
