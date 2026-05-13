package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration for the service.
// Values are loaded from environment variables (with an optional .env file fallback).
type Config struct {
	App      AppConfig
	HTTP     HTTPConfig
	Database DatabaseConfig
	Log      LogConfig
}

// AppConfig holds application-level settings.
type AppConfig struct {
	Name  string `mapstructure:"APP_NAME"`
	Env   string `mapstructure:"APP_ENV"`   // local | development | staging | production
	Debug bool   `mapstructure:"APP_DEBUG"`
}

// LogConfig holds logger settings.
type LogConfig struct {
	// Level is the minimum log level to emit: debug | info | warn | error
	Level string `mapstructure:"LOG_LEVEL"`
}

// HTTPConfig holds HTTP server settings.
type HTTPConfig struct {
	Host           string `mapstructure:"HTTP_HOST"`
	Port           string `mapstructure:"HTTP_PORT"`
	AllowedOrigins string `mapstructure:"HTTP_ALLOWED_ORIGINS"`
}

// DatabaseConfig holds database connection settings.
type DatabaseConfig struct {
	DSN string `mapstructure:"DATABASE_DSN"`
}

// Addr returns the full host:port address for the HTTP server.
func (c HTTPConfig) Addr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

// Load reads configuration from environment variables and an optional .env file.
// Environment variables always take precedence over file values.
func Load() (*Config, error) {
	v := viper.New()

	// --- defaults ---
	v.SetDefault("APP_NAME", "product-service")
	v.SetDefault("APP_ENV", "local")
	v.SetDefault("APP_DEBUG", false)
	v.SetDefault("HTTP_HOST", "")
	v.SetDefault("HTTP_PORT", "8080")
	v.SetDefault("HTTP_ALLOWED_ORIGINS", "*")
	v.SetDefault("DATABASE_DSN", "")
	v.SetDefault("LOG_LEVEL", "info")

	// --- .env file (optional, does not fail if missing) ---
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	_ = v.ReadInConfig() // intentionally ignore "file not found" errors

	// --- environment variables override everything ---
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("config: failed to unmarshal: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate checks that required fields are present.
func (c *Config) validate() error {
	if c.HTTP.Port == "" {
		return fmt.Errorf("config: HTTP_PORT must not be empty")
	}
	return nil
}
