package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user account
type User struct {
	BaseModel

	Email        string `gorm:"type:varchar(255);uniqueIndex;not null"`
	Username     string `gorm:"type:varchar(255);uniqueIndex"`
	PasswordHash string `gorm:"type:varchar(255);not null"`
	Name         string `gorm:"type:varchar(255)"`
	IsActive     bool   `gorm:"default:true"`
	IsAdmin      bool   `gorm:"default:false"`
	LastLoginAt  *time.Time

	// Associations
	Links   []Link   `gorm:"foreignKey:UserID"`
	Domains []Domain `gorm:"foreignKey:UserID"`
	APIKeys []APIKey `gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for GORM
func (User) TableName() string {
	return "users"
}

// BeforeCreate validates the user before creation
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// Call BaseModel's BeforeCreate to set UUID
	if err := u.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// Ensure email is not empty
	if u.Email == "" {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// IsVerified checks if the user is verified (placeholder for future email verification)
func (u *User) IsVerified() bool {
	// For now, all users are considered verified
	// In the future, add email verification field
	return true
}

// CanCreateLink checks if the user can create a new link
func (u *User) CanCreateLink(currentCount int, maxLimit int) bool {
	if !u.IsActive {
		return false
	}

	if maxLimit > 0 && currentCount >= maxLimit {
		return false
	}

	return true
}
