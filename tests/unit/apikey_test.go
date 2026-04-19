package unit

import (
	"testing"
	"time"

	"url-shortener/internal/models"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestAPIKey_IsExpired(t *testing.T) {
	now := time.Now()
	past := now.Add(-1 * time.Hour)
	future := now.Add(1 * time.Hour)

	tests := []struct {
		name     string
		expires  *time.Time
		expected bool
	}{
		{"No expiry", nil, false},
		{"Expired", &past, true},
		{"Not expired", &future, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := models.APIKey{
				ExpiresAt: tt.expires,
			}
			assert.Equal(t, tt.expected, key.IsExpired())
		})
	}
}

func TestAPIKey_CanBeUsed(t *testing.T) {
	now := time.Now()
	future := now.Add(1 * time.Hour)
	past := now.Add(-1 * time.Hour)

	key := models.APIKey{
		ExpiresAt: &future,
	}
	assert.True(t, key.CanBeUsed())

	// Expired key
	key.ExpiresAt = &past
	assert.False(t, key.CanBeUsed())

	// No expiry (nil)
	key.ExpiresAt = nil
	assert.True(t, key.CanBeUsed())
}

func TestAPIKey_UpdateLastUsed(t *testing.T) {
	key := models.APIKey{}
	assert.Nil(t, key.LastUsedAt)

	key.UpdateLastUsed()
	assert.NotNil(t, key.LastUsedAt)
	assert.WithinDuration(t, time.Now(), *key.LastUsedAt, time.Second)
}

func TestAPIKey_BeforeCreate(t *testing.T) {
	db := &gorm.DB{}
	key := models.APIKey{
		KeyHash: "hash",
		Name:    "Test Key",
	}
	err := key.BeforeCreate(db)
	assert.NoError(t, err)

	// Missing key hash
	key2 := models.APIKey{Name: "Test"}
	err = key2.BeforeCreate(db)
	assert.Error(t, err)

	// Missing name
	key3 := models.APIKey{KeyHash: "hash"}
	err = key3.BeforeCreate(db)
	assert.Error(t, err)
}

func TestAPIKey_TableName(t *testing.T) {
	key := models.APIKey{}
	assert.Equal(t, "api_keys", key.TableName())
}