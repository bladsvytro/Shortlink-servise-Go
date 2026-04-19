package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// APIKey represents an API key for programmatic access
type APIKey struct {
	BaseModel

	UserID     uuid.UUID `gorm:"type:uuid;index;not null"`
	KeyHash    string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	Name       string    `gorm:"type:varchar(255);not null"`
	LastUsedAt *time.Time
	ExpiresAt  *time.Time
	RateLimit  int `gorm:"default:1000"`

	// Associations
	User User `gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for GORM
func (APIKey) TableName() string {
	return "api_keys"
}

// BeforeCreate validates the API key before creation
func (a *APIKey) BeforeCreate(tx *gorm.DB) error {
	// Call BaseModel's BeforeCreate to set UUID
	if err := a.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// Ensure key hash and name are not empty
	if a.KeyHash == "" || a.Name == "" {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// IsExpired checks if the API key has expired
func (a *APIKey) IsExpired() bool {
	if a.ExpiresAt == nil {
		return false
	}
	return a.ExpiresAt.Before(time.Now())
}

// CanBeUsed checks if the API key can be used
func (a *APIKey) CanBeUsed() bool {
	return !a.IsExpired()
}

// UpdateLastUsed updates the LastUsedAt timestamp
func (a *APIKey) UpdateLastUsed() {
	now := time.Now()
	a.LastUsedAt = &now
}
