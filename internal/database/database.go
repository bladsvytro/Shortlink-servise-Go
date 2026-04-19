package database

import (
	"fmt"
	"time"

	"url-shortener/internal/config"
	"url-shortener/internal/models"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Database represents a database connection
type Database struct {
	DB     *gorm.DB
	config config.DatabaseConfig
	logger *zap.Logger
}

// New creates a new database connection
func New(cfg config.DatabaseConfig, logger *zap.Logger) (*Database, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL DB for connection pooling
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying SQL DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(cfg.MaxConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	database := &Database{
		DB:     db,
		config: cfg,
		logger: logger,
	}

	// Test connection
	if err := database.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	logger.Info("Database connected",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.String("database", cfg.Name),
	)

	return database, nil
}

// Ping checks if the database is reachable
func (d *Database) Ping() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// Close closes the database connection
func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Migrate runs database migrations
func (d *Database) Migrate() error {
	d.logger.Info("Running database migrations...")

	// Auto-migrate all models
	err := d.DB.AutoMigrate(
		&models.User{},
		&models.Link{},
		&models.Domain{},
		&models.APIKey{},
	)
	if err != nil {
		d.logger.Error("Migration failed", zap.Error(err))
		return fmt.Errorf("migration failed: %w", err)
	}

	// Create anonymous user if not exists
	anonymousID := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	var anonymousUser models.User
	if err := d.DB.First(&anonymousUser, "id = ?", anonymousID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			d.logger.Info("Creating anonymous user...")
			anonymousUser = models.User{
				BaseModel: models.BaseModel{
					ID: anonymousID,
				},
				Email:        "anonymous@example.com",
				PasswordHash: "anonymous", // Not used for authentication
				Name:         "Anonymous User",
				IsActive:     true,
				IsAdmin:      false,
			}
			// Create with skipped hooks to preserve the zero UUID
			if err := d.DB.Session(&gorm.Session{SkipHooks: true}).Create(&anonymousUser).Error; err != nil {
				d.logger.Error("Failed to create anonymous user", zap.Error(err))
				return fmt.Errorf("failed to create anonymous user: %w", err)
			}
			d.logger.Info("Anonymous user created")
		} else {
			d.logger.Error("Failed to query anonymous user", zap.Error(err))
			return fmt.Errorf("failed to query anonymous user: %w", err)
		}
	}

	d.logger.Info("Database migrations completed")
	return nil
}

// Transaction executes a function within a database transaction
func (d *Database) Transaction(fn func(tx *gorm.DB) error) error {
	return d.DB.Transaction(fn)
}
