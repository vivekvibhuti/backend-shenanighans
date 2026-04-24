// Package config provides application configuration management.
//
// It loads configuration from environment variables and YAML config files,
// with support for defaults and validation. Configuration is loaded once
// at application startup.
// not using functional/optional pattern for this,
// just load the config with Load(), and it returns a struct
//
// # Configuration Precedence (highest to lowest)
//
//   - Environment variables (APP_* prefix)
//   - Config file (config.yaml)
//   - Defaults
//
// # Usage
//
//	cfg, err := config.Load()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	cfg.Database.URL // database connection string
//	cfg.Server.Port  // e.g. "8080"
package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

var (
	ErrMissingDatabaseURL = errors.New("database URL is required")
	ErrMissingJWTSecret   = errors.New("JWT secret is required")
	ErrJWTSecretTooShort  = errors.New("JWT secret must be at least 32 characters")
)

// Config holds all application configuration.
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Logger   LoggerConfig   `mapstructure:"logger"`
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Port         string        `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// DatabaseConfig holds database connection settings.
type DatabaseConfig struct {
	URL          string `mapstructure:"url"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

// JWTConfig holds JWT authentication settings.
type JWTConfig struct {
	Secret         string        `mapstructure:"secret"`
	ExpiryDuration time.Duration `mapstructure:"expiry_duration"`
}

// LoggerConfig holds logging settings.
type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

// Load reads and validates configuration from all sources.
// It returns an error if required configuration is missing or invalid.
func Load() (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")

	v.SetDefault("server.port", "8080")
	v.SetDefault("server.read_timeout", 15*time.Second)
	v.SetDefault("server.write_timeout", 15*time.Second)
	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 5)
	v.SetDefault("jwt.expiry_duration", 24*time.Hour)
	v.SetDefault("logger.level", "info")

	v.SetEnvPrefix("APP")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("read config: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// validate checks that required configuration is present and valid.
func (c *Config) validate() error {
	if c.Database.URL == "" {
		return ErrMissingDatabaseURL
	}
	if c.JWT.Secret == "" {
		return ErrMissingJWTSecret
	}
	if len(c.JWT.Secret) < 32 {
		return ErrJWTSecretTooShort
	}
	return nil
}
