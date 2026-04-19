package unit

import (
	"testing"
	"time"

	"url-shortener/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUserModel(t *testing.T) {
	user := models.User{
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Name:         "Test User",
		IsActive:     true,
		IsAdmin:      false,
	}
	assert.Equal(t, "users", user.TableName())
	assert.True(t, user.IsVerified())
	assert.True(t, user.CanCreateLink(0, 10))
	assert.False(t, user.CanCreateLink(10, 10))
	assert.False(t, user.CanCreateLink(11, 10))
	user.IsActive = false
	assert.False(t, user.CanCreateLink(0, 10))
}

func TestLinkModel(t *testing.T) {
	link := models.Link{
		ShortCode:   "abc123",
		OriginalURL: "https://example.com",
		UserID:      uuid.New(),
		Title:       "Example",
		Description: "Test link",
		IsActive:    true,
		ClickCount:  0,
	}
	assert.Equal(t, "links", link.TableName())

	// Test GetShortURL
	link.ShortCode = "test"
	assert.Equal(t, "http://localhost:8080/test", link.GetShortURL("http://localhost:8080"))
	assert.Equal(t, "https://short.io/test", link.GetShortURL("https://short.io"))
}

func TestBaseModel(t *testing.T) {
	base := models.BaseModel{}
	// Initially zero values
	assert.Equal(t, uuid.Nil, base.ID)
	assert.Zero(t, base.CreatedAt)
	assert.Zero(t, base.UpdatedAt)

	// Simulate BeforeCreate
	db := &gorm.DB{}
	err := base.BeforeCreate(db)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, base.ID)
	// CreatedAt and UpdatedAt are set by GORM automatically, but we don't have a real DB.
	// We'll just ensure ID is set.
}

func TestDomainModel(t *testing.T) {
	domain := models.Domain{
		DomainName: "example.com",
		UserID:     uuid.New(),
		IsActive:   true,
	}
	assert.Equal(t, "domains", domain.TableName())
}

func TestAPIKeyModel(t *testing.T) {
	expiresAt := time.Now().Add(24 * time.Hour)
	apiKey := models.APIKey{
		KeyHash:   "test_key_123",
		UserID:    uuid.New(),
		Name:      "Test Key",
		ExpiresAt: &expiresAt,
	}
	assert.Equal(t, "api_keys", apiKey.TableName())
}
