package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseModel contains common fields for all models
type BaseModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (base *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if base.ID == uuid.Nil {
		base.ID = uuid.New()
	}
	return nil
}

// Timestamps interface for models that need timestamps
type Timestamps interface {
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
}

// GetCreatedAt returns the CreatedAt time
func (base *BaseModel) GetCreatedAt() time.Time {
	return base.CreatedAt
}

// GetUpdatedAt returns the UpdatedAt time
func (base *BaseModel) GetUpdatedAt() time.Time {
	return base.UpdatedAt
}
