package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"url-shortener/internal/app"
	"url-shortener/internal/config"
	"url-shortener/internal/models"
	"url-shortener/internal/pkg/logger"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestApp(t *testing.T) (*app.Application, func()) {
	// Create a test configuration
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host:         "localhost",
			Port:         8080,
			Env:          "test",
			ReadTimeout:  5,
			WriteTimeout: 10,
			IdleTimeout:  60,
		},
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "testuser",
			Password: "testpass",
			Name:     "testdb",
			SSLMode:  "disable",
		},
		Auth: config.AuthConfig{
			JWTSecret:          "test-jwt-secret-key-for-integration-tests",
			AccessTokenExpiry:  3600,
			RefreshTokenExpiry: 86400,
			BCryptCost:         4, // lower cost for faster tests
		},
		Redis: config.RedisConfig{
			Enabled: false,
		},
	}

	// Create a logger
	log, err := logger.New("debug", "console")
	require.NoError(t, err)

	// Create application
	app, err := app.New(cfg, log)
	require.NoError(t, err)

	// Run migrations
	err = app.DB().Migrate()
	require.NoError(t, err)

	// Cleanup function
	cleanup := func() {
		// Close database connection
		app.DB().Close()
	}

	return app, cleanup
}

func TestHealthEndpoint(t *testing.T) {
	app, cleanup := setupTestApp(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	app.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "OK", w.Body.String())
}

func TestRedirectEndpoint(t *testing.T) {
	app, cleanup := setupTestApp(t)
	defer cleanup()

	// Create a link directly in the database
	link := &models.Link{
		ShortCode:   "test123",
		OriginalURL: "https://example.com",
		UserID:      uuid.Nil, // anonymous
		ClickCount:  0,
		ExpiresAt:   nil,
	}
	err := app.DB().DB.Create(link).Error
	require.NoError(t, err)

	// Test redirect
	req := httptest.NewRequest("GET", "/test123", nil)
	w := httptest.NewRecorder()
	app.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Contains(t, w.Header().Get("Location"), "https://example.com")

	// Verify click count increased
	var updatedLink models.Link
	err = app.DB().DB.Where("short_code = ?", "test123").First(&updatedLink).Error
	require.NoError(t, err)
	assert.Equal(t, int64(1), updatedLink.ClickCount)
}

func TestRedirectEndpoint_NotFound(t *testing.T) {
	app, cleanup := setupTestApp(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/nonexistent", nil)
	w := httptest.NewRecorder()
	app.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAuthRegisterAndLogin(t *testing.T) {
	app, cleanup := setupTestApp(t)
	defer cleanup()

	// Test registration
	registerBody := `{"email":"test@example.com","password":"password123"}`
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBufferString(registerBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	app.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Test login
	loginBody := `{"email":"test@example.com","password":"password123"}`
	req = httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBufferString(loginBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	app.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var loginResp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &loginResp)
	require.NoError(t, err)
	assert.NotEmpty(t, loginResp["access_token"])
}