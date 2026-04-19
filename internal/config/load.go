package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Load loads configuration from file and environment variables
func Load() (*Config, error) {
	// Load .env file if exists
	_ = godotenv.Load(".env.local")
	_ = godotenv.Load()

	// Setup Viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Read config file (optional)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Set defaults first
	setDefaults()

	// Bind environment variables
	viper.AutomaticEnv()
	bindEnvVars()

	// Create config with defaults
	cfg := DefaultConfig()

	// Unmarshal config
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Validate config
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// setDefaults sets Viper defaults
func setDefaults() {
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.env", "development")
	viper.SetDefault("server.read_timeout", 30)
	viper.SetDefault("server.write_timeout", 30)
	viper.SetDefault("server.idle_timeout", 120)

	viper.SetDefault("database.host", "postgres")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "urlshortener")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.name", "url_shortener")
	viper.SetDefault("database.ssl_mode", "disable")
	viper.SetDefault("database.max_conns", 25)

	viper.SetDefault("auth.jwt_secret", "your-super-secret-jwt-key-change-in-production")
	viper.SetDefault("auth.access_token_expiry", 900)
	viper.SetDefault("auth.refresh_token_expiry", 604800)
	viper.SetDefault("auth.bcrypt_cost", 12)

	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "console")

	viper.SetDefault("shortener.base_url", "http://localhost:8080")
	viper.SetDefault("shortener.short_code_length", 6)
	viper.SetDefault("shortener.allow_custom_codes", true)
	viper.SetDefault("shortener.max_custom_code_length", 32)
	viper.SetDefault("shortener.default_redirect_code", 302)
}

// bindEnvVars explicitly binds environment variables to Viper keys
func bindEnvVars() {
	viper.BindEnv("server.port", "SERVER_PORT")
	viper.BindEnv("server.host", "SERVER_HOST")
	viper.BindEnv("server.env", "SERVER_ENV")
	viper.BindEnv("server.read_timeout", "SERVER_READ_TIMEOUT")
	viper.BindEnv("server.write_timeout", "SERVER_WRITE_TIMEOUT")
	viper.BindEnv("server.idle_timeout", "SERVER_IDLE_TIMEOUT")

	viper.BindEnv("database.host", "DATABASE_HOST")
	viper.BindEnv("database.port", "DATABASE_PORT")
	viper.BindEnv("database.user", "DATABASE_USER")
	viper.BindEnv("database.password", "DATABASE_PASSWORD")
	viper.BindEnv("database.name", "DATABASE_NAME")
	viper.BindEnv("database.ssl_mode", "DATABASE_SSL_MODE")
	viper.BindEnv("database.max_conns", "DATABASE_MAX_CONNS")

	viper.BindEnv("auth.jwt_secret", "AUTH_JWT_SECRET")
	viper.BindEnv("auth.access_token_expiry", "AUTH_ACCESS_TOKEN_EXPIRY")
	viper.BindEnv("auth.refresh_token_expiry", "AUTH_REFRESH_TOKEN_EXPIRY")
	viper.BindEnv("auth.bcrypt_cost", "AUTH_BCRYPT_COST")

	viper.BindEnv("logging.level", "LOGGING_LEVEL")
	viper.BindEnv("logging.format", "LOGGING_FORMAT")

	viper.BindEnv("shortener.base_url", "SHORTENER_BASE_URL")
	viper.BindEnv("shortener.short_code_length", "SHORTENER_SHORT_CODE_LENGTH")
	viper.BindEnv("shortener.allow_custom_codes", "SHORTENER_ALLOW_CUSTOM_CODES")
	viper.BindEnv("shortener.max_custom_code_length", "SHORTENER_MAX_CUSTOM_CODE_LENGTH")
	viper.BindEnv("shortener.default_redirect_code", "SHORTENER_DEFAULT_REDIRECT_CODE")
}
