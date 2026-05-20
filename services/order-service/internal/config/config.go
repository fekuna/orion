package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	HTTP     HTTPConfig
	Database DatabaseConfig
}

type AppConfig struct {
	Name  string `mapstructure:"APP_NAME"`
	Env   string `mapstructure:"APP_ENV"`
	Debug bool   `mapstructure:"APP_DEBUG"`
}

type HTTPConfig struct {
	Host           string `mapstructure:"HTTP_HOST"`
	Port           string `mapstructure:"HTTP_PORT"`
	AllowedOrigins string `mapstructure:"HTTP_ALLOWED_ORIGINS"`
}

type LogConfig struct {
	Level string `mapstructure:"LOG_LEVEL"`
}

type DatabaseConfig struct {
	DSN string `mapstructure:"DATABASE_DSN"`
}

func (c HTTPConfig) Addr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func Load() (*Config, error) {
	v := viper.New()
	// --- defaults ---
	v.SetDefault("APP_NAME", "order-service")
	v.SetDefault("APP_ENV", "local")
	v.SetDefault("APP_DEBUG", false)
	v.SetDefault("HTTP_HOST", "")
	v.SetDefault("HTTP_PORT", "8080")
	v.SetDefault("HTTP_ALLOWED_ORIGINS", "*")
	v.SetDefault("DATABASE_DSN", "")
	v.SetDefault("LOG_LEVEL", "info")

	// --- env file ---
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	_ = v.ReadInConfig()

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
