package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Link represents a shortened URL
type Link struct {
	BaseModel

	ShortCode     string     `gorm:"type:varchar(32);uniqueIndex;not null"`
	OriginalURL   string     `gorm:"type:text;not null"`
	UserID        uuid.UUID  `gorm:"type:uuid;index"`
	DomainID      *uuid.UUID `gorm:"type:uuid;index"`
	Title         string     `gorm:"type:varchar(255)"`
	Description   string     `gorm:"type:text"`
	Tags          []string   `gorm:"type:text[]"`
	IsActive      bool       `gorm:"default:true"`
	ExpiresAt     *time.Time
	ClickCount    int64 `gorm:"default:0"`
	LastClickedAt *time.Time

	// Associations (will be populated by GORM)
	User   User    `gorm:"foreignKey:UserID"`
	Domain *Domain `gorm:"foreignKey:DomainID"`
}

// TableName specifies the table name for GORM
func (Link) TableName() string {
	return "links"
}

// BeforeCreate validates the link before creation
func (l *Link) BeforeCreate(tx *gorm.DB) error {
	// Call BaseModel's BeforeCreate to set UUID
	if err := l.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// Ensure short code is not empty
	if l.ShortCode == "" {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// IsExpired checks if the link has expired
func (l *Link) IsExpired() bool {
	if l.ExpiresAt == nil {
		return false
	}
	return l.ExpiresAt.Before(time.Now())
}

// CanBeAccessed checks if the link can be accessed (active and not expired)
func (l *Link) CanBeAccessed() bool {
	return l.IsActive && !l.IsExpired()
}

// IncrementClickCount increments the click count and updates LastClickedAt
func (l *Link) IncrementClickCount() {
	l.ClickCount++
	now := time.Now()
	l.LastClickedAt = &now
}

// GetShortURL returns the full short URL
func (l *Link) GetShortURL(baseURL string) string {
	if l.Domain != nil && l.Domain.DomainName != "" {
		return "https://" + l.Domain.DomainName + "/" + l.ShortCode
	}
	return baseURL + "/" + l.ShortCode
}
