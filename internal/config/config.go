package config

import (
	"fmt"
)

// Config represents the application configuration
type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Auth      AuthConfig      `mapstructure:"auth"`
	Redis     RedisConfig     `mapstructure:"redis"`
	Logging   LoggingConfig   `mapstructure:"logging"`
	Shortener ShortenerConfig `mapstructure:"shortener"`
}

// ServerConfig contains HTTP server configuration
type ServerConfig struct {
	Port         int    `mapstructure:"port"`
	Host         string `mapstructure:"host"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
	IdleTimeout  int    `mapstructure:"idle_timeout"`
	Env          string `mapstructure:"env"`
}

// DatabaseConfig contains PostgreSQL database configuration
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"ssl_mode"`
	MaxConns int    `mapstructure:"max_conns"`
}

// AuthConfig contains authentication configuration
type AuthConfig struct {
	JWTSecret          string `mapstructure:"jwt_secret"`
	AccessTokenExpiry  int    `mapstructure:"access_token_expiry"`
	RefreshTokenExpiry int    `mapstructure:"refresh_token_expiry"`
	BCryptCost         int    `mapstructure:"bcrypt_cost"`
}

// RedisConfig contains Redis configuration (optional)
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	Enabled  bool   `mapstructure:"enabled"`
}

// LoggingConfig contains logging configuration
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"` // json or console
}

// ShortenerConfig contains URL shortener specific configuration
type ShortenerConfig struct {
	BaseURL             string `mapstructure:"base_url"`
	ShortCodeLength     int    `mapstructure:"short_code_length"`
	AllowCustomCodes    bool   `mapstructure:"allow_custom_codes"`
	MaxCustomCodeLength int    `mapstructure:"max_custom_code_length"`
	DefaultRedirectCode int    `mapstructure:"default_redirect_code"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         8080,
			Host:         "0.0.0.0",
			ReadTimeout:  30,
			WriteTimeout: 30,
			IdleTimeout:  120,
			Env:          "development",
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "urlshortener",
			Password: "password",
			Name:     "url_shortener",
			SSLMode:  "disable",
			MaxConns: 25,
		},
		Auth: AuthConfig{
			JWTSecret:          "your-super-secret-jwt-key-change-in-production",
			AccessTokenExpiry:  900,
			RefreshTokenExpiry: 604800,
			BCryptCost:         12,
		},
		Redis: RedisConfig{
			Host:     "localhost",
			Port:     6379,
			Password: "",
			DB:       0,
			Enabled:  false,
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "console",
		},
		Shortener: ShortenerConfig{
			BaseURL:             "http://localhost:8080",
			ShortCodeLength:     6,
			AllowCustomCodes:    true,
			MaxCustomCodeLength: 32,
			DefaultRedirectCode: 302,
		},
	}
}

// Validate performs basic validation of configuration
func (c *Config) Validate() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}
	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if c.Database.Name == "" {
		return fmt.Errorf("database name is required")
	}
	if c.Auth.JWTSecret == "" {
		return fmt.Errorf("JWT secret is required")
	}
	if c.Shortener.BaseURL == "" {
		return fmt.Errorf("shortener base URL is required")
	}
	return nil
}
