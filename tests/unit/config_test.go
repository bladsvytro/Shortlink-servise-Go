package unit

import (
	"os"
	"testing"

	"url-shortener/internal/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	// Temporarily set some environment variables to override defaults
	t.Setenv("SERVER_PORT", "9090")
	t.Setenv("SERVER_HOST", "127.0.0.1")
	t.Setenv("DATABASE_HOST", "testhost")
	t.Setenv("AUTH_JWT_SECRET", "test-jwt-secret")

	// Load config
	cfg, err := config.Load()
	require.NoError(t, err)
	assert.NotNil(t, cfg)

	// Verify overridden values
	assert.Equal(t, 9090, cfg.Server.Port)
	assert.Equal(t, "127.0.0.1", cfg.Server.Host)
	assert.Equal(t, "testhost", cfg.Database.Host)
	assert.Equal(t, "test-jwt-secret", cfg.Auth.JWTSecret)

	// Reset environment after test
	t.Cleanup(func() {
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("SERVER_HOST")
		os.Unsetenv("DATABASE_HOST")
		os.Unsetenv("AUTH_JWT_SECRET")
	})
}

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	assert.NotNil(t, cfg)
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "0.0.0.0", cfg.Server.Host)
	assert.Equal(t, "development", cfg.Server.Env)
	assert.Equal(t, "localhost", cfg.Database.Host)
	assert.Equal(t, 5432, cfg.Database.Port)
	assert.Equal(t, "urlshortener", cfg.Database.User)
	assert.Equal(t, "password", cfg.Database.Password)
	assert.Equal(t, "url_shortener", cfg.Database.Name)
}

func TestConfigValidation(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port: 0, // invalid
		},
	}
	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "server port")

	cfg.Server.Port = 8080
	cfg.Database.Host = ""
	err = cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database host")
}
