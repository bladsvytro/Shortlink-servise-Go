package unit

import (
	"testing"
	"time"

	"url-shortener/internal/app"
	"url-shortener/internal/config"
	"url-shortener/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	password := "mysecretpassword"
	hash, err := app.HashPassword(password)
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.True(t, app.CheckPasswordHash(password, hash))
	assert.False(t, app.CheckPasswordHash("wrongpassword", hash))
}

func TestGenerateAndValidateJWT(t *testing.T) {
	cfg := config.AuthConfig{
		JWTSecret:          "test-secret-key",
		AccessTokenExpiry:  3600,
		RefreshTokenExpiry: 86400,
		BCryptCost:         12,
	}
	user := &models.User{
		Email:   "test@example.com",
		IsAdmin: false,
	}
	user.ID = uuid.New()

	token, err := app.GenerateJWT(user, &cfg)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := app.ValidateJWT(token, &cfg)
	require.NoError(t, err)
	assert.Equal(t, user.ID.String(), claims.UserID)
	assert.Equal(t, user.Email, claims.Email)
	assert.Equal(t, user.IsAdmin, claims.IsAdmin)
	assert.WithinDuration(t, time.Now().Add(time.Duration(cfg.AccessTokenExpiry)*time.Second), claims.ExpiresAt.Time, 5*time.Second)
}

func TestValidateJWT_Invalid(t *testing.T) {
	cfg := config.AuthConfig{
		JWTSecret: "test-secret-key",
	}
	// Invalid token
	_, err := app.ValidateJWT("invalid.token.here", &cfg)
	assert.Error(t, err)

	// Token with wrong secret
	user := &models.User{
		Email: "test@example.com",
	}
	user.ID = uuid.New()
	token, _ := app.GenerateJWT(user, &cfg)
	cfg2 := config.AuthConfig{JWTSecret: "different-secret"}
	_, err = app.ValidateJWT(token, &cfg2)
	assert.Error(t, err)
}

func TestAuthClaims(t *testing.T) {
	claims := &app.AuthClaims{
		UserID:  "123e4567-e89b-12d3-a456-426614174000",
		Email:   "user@example.com",
		IsAdmin: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			Issuer:    "test",
		},
	}
	assert.Equal(t, "123e4567-e89b-12d3-a456-426614174000", claims.UserID)
	assert.Equal(t, "user@example.com", claims.Email)
	assert.True(t, claims.IsAdmin)
}
