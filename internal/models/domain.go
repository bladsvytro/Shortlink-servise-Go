package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Domain represents a custom domain for URL shortening
type Domain struct {
	BaseModel

	DomainName        string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	UserID            uuid.UUID `gorm:"type:uuid;index;not null"`
	IsVerified        bool      `gorm:"default:false"`
	IsActive          bool      `gorm:"default:true"`
	VerificationToken string    `gorm:"type:varchar(64)"`
	VerifiedAt        *time.Time

	// Associations
	User  User   `gorm:"foreignKey:UserID"`
	Links []Link `gorm:"foreignKey:DomainID"`
}

// TableName specifies the table name for GORM
func (Domain) TableName() string {
	return "domains"
}

// BeforeCreate validates the domain before creation
func (d *Domain) BeforeCreate(tx *gorm.DB) error {
	// Call BaseModel's BeforeCreate to set UUID
	if err := d.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// Ensure domain name is not empty
	if d.DomainName == "" {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// CanBeUsed checks if the domain can be used for shortening
func (d *Domain) CanBeUsed() bool {
	return d.IsVerified && d.IsActive
}

// Verify marks the domain as verified
func (d *Domain) Verify() {
	d.IsVerified = true
	now := time.Now()
	d.VerifiedAt = &now
}
