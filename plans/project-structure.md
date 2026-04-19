# Project Structure for Go URL Shortener

## Complete Directory Tree

```
url-shortener/
├── cmd/
│   └── server/
│       ├── main.go                 # Application entry point
│       └── wire.go                 # Dependency injection (optional)
├── internal/
│   ├── config/
│   │   ├── config.go              # Configuration structs and loading
│   │   ├── load.go                # Viper configuration loader
│   │   └── env.go                 # Environment variable validation
│   ├── database/
│   │   ├── connection.go          # Database connection setup
│   │   ├── migrations/            # Auto-generated migration files
│   │   │   ├── 001_init_schema.up.sql
│   │   │   └── 001_init_schema.down.sql
│   │   └── seed/                  # Seed data for development
│   │       └── seed.go
│   ├── models/
│   │   ├── user.go                # User model
│   │   ├── link.go                # Link model
│   │   ├── domain.go              # Domain model
│   │   ├── api_key.go             # API Key model
│   │   ├── click_event.go         # Click Event model
│   │   └── base.go                # Base model with common fields
│   ├── repositories/
│   │   ├── interface.go           # Repository interfaces
│   │   ├── user_repository.go     # User data access
│   │   ├── link_repository.go     # Link data access
│   │   ├── domain_repository.go   # Domain data access
│   │   ├── analytics_repository.go # Analytics data access
│   │   └── base_repository.go     # Base repository with common methods
│   ├── services/
│   │   ├── auth_service.go        # Authentication logic
│   │   ├── link_service.go        # Link generation and management
│   │   ├── domain_service.go      # Domain verification and management
│   │   ├── analytics_service.go   # Analytics processing
│   │   ├── shortener_service.go   # URL shortening algorithm
│   │   ├── validator_service.go   # Input validation
│   │   └── cache_service.go       # Caching layer (optional)
│   ├── handlers/
│   │   ├── middleware/
│   │   │   ├── auth.go            # JWT authentication middleware
│   │   │   ├── logging.go         # Request logging
│   │   │   ├── rate_limit.go      # Rate limiting
│   │   │   ├── recovery.go        # Panic recovery
│   │   │   └── cors.go            # CORS configuration
│   │   ├── auth_handler.go        # Authentication endpoints
│   │   ├── link_handler.go        # Link management endpoints
│   │   ├── domain_handler.go      # Domain management endpoints
│   │   ├── analytics_handler.go   # Analytics endpoints
│   │   ├── redirect_handler.go    # Public redirect endpoint
│   │   ├── admin_handler.go       # Admin endpoints
│   │   └── health_handler.go      # Health check endpoints
│   ├── utils/
│   │   ├── jwt.go                 # JWT token generation/validation
│   │   ├── password.go            # Password hashing and verification
│   │   ├── validator.go           # Custom validation rules
│   │   ├── shortcode.go           # Short code generation
│   │   ├── url.go                 # URL validation and normalization
│   │   ├── pagination.go          # Pagination utilities
│   │   ├── response.go            # Standardized API responses
│   │   ├── error.go               # Error handling utilities
│   │   └── time.go                # Time utilities
│   ├── pkg/                       # Reusable packages (if any)
│   │   ├── logger/
│   │   │   └── logger.go          # Structured logger setup
│   │   └── metrics/
│   │       └── metrics.go         # Prometheus metrics setup
│   └── app/                       # Application layer
│       ├── app.go                 # Application struct and lifecycle
│       └── server.go              # HTTP server setup
├── api/
│   └── v1/
│       ├── openapi.yaml           # OpenAPI 3.0 specification
│       ├── openapi.json           # JSON version of OpenAPI spec
│       └── client/                # Generated API client (optional)
│           ├── go.mod
│           └── client.go
├── migrations/                    # Manual migration files
│   ├── 001_init_schema.up.sql
│   ├── 001_init_schema.down.sql
│   ├── 002_add_analytics.up.sql
│   └── 002_add_analytics.down.sql
├── tests/
│   ├── unit/
│   │   ├── services/
│   │   │   ├── auth_service_test.go
│   │   │   └── link_service_test.go
│   │   ├── utils/
│   │   │   └── shortcode_test.go
│   │   └── handlers/
│   │       └── auth_handler_test.go
│   ├── integration/
│   │   ├── database_test.go
│   │   ├── api_test.go
│   │   └── test_helpers.go
│   ├── e2e/
│   │   └── api_e2e_test.go
│   └── fixtures/                  # Test data
│       ├── users.json
│       └── links.json
├── docker/
│   ├── Dockerfile                 # Production Dockerfile
│   ├── Dockerfile.dev             # Development Dockerfile
│   ├── docker-compose.yml         # Local development stack
│   └── docker-compose.test.yml    # Test environment
├── deployments/
│   ├── kubernetes/
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   ├── ingress.yaml
│   │   └── configmap.yaml
│   └── helm/
│       ├── Chart.yaml
│       ├── values.yaml
│       └── templates/
├── scripts/
│   ├── migrate.sh                 # Database migration script
│   ├── seed.sh                    # Seed database script
│   ├── test.sh                    # Test runner script
│   ├── build.sh                   # Build script
│   └── deploy.sh                  # Deployment script
├── web/                           # Frontend (optional)
│   ├── public/
│   └── src/
├── docs/
│   ├── api.md                     # API documentation
│   ├── architecture.md            # Architecture decisions
│   ├── deployment.md              # Deployment guide
│   └── development.md             # Development guide
├── .env.example                   # Environment variables template
├── .env.local                     # Local development env (gitignored)
├── .gitignore
├── .golangci.yml                  # Linter configuration
├── .air.toml                      # Hot reload configuration
├── Makefile                       # Task automation
├── go.mod                         # Go module definition
├── go.sum
├── README.md                      # Project overview
└── CHANGELOG.md                   # Version history
```

## Key File Examples

### `cmd/server/main.go`
```go
package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"
    
    "github.com/your-org/url-shortener/internal/app"
    "github.com/your-org/url-shortener/internal/config"
    "github.com/your-org/url-shortener/pkg/logger"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    // Initialize logger
    zapLogger, err := logger.New(cfg.LogLevel)
    if err != nil {
        log.Fatalf("Failed to initialize logger: %v", err)
    }
    defer zapLogger.Sync()
    
    // Create application
    application, err := app.New(cfg, zapLogger)
    if err != nil {
        zapLogger.Fatal("Failed to create application", zap.Error(err))
    }
    
    // Start application
    if err := application.Start(); err != nil {
        zapLogger.Fatal("Failed to start application", zap.Error(err))
    }
    
    // Graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    zapLogger.Info("Shutting down server...")
    application.Stop()
    zapLogger.Info("Server stopped")
}
```

### `internal/config/config.go`
```go
package config

type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Auth     AuthConfig
    Redis    RedisConfig `mapstructure:"redis"`
    Logging  LoggingConfig
}

type ServerConfig struct {
    Port         int    `mapstructure:"port"`
    Host         string `mapstructure:"host"`
    ReadTimeout  int    `mapstructure:"read_timeout"`
    WriteTimeout int    `mapstructure:"write_timeout"`
    IdleTimeout  int    `mapstructure:"idle_timeout"`
    Env          string `mapstructure:"env"`
}

type DatabaseConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    User     string `mapstructure:"user"`
    Password string `mapstructure:"password"`
    Name     string `mapstructure:"name"`
    SSLMode  string `mapstructure:"ssl_mode"`
    MaxConns int    `mapstructure:"max_conns"`
}

type AuthConfig struct {
    JWTSecret          string `mapstructure:"jwt_secret"`
    AccessTokenExpiry  int    `mapstructure:"access_token_expiry"`
    RefreshTokenExpiry int    `mapstructure:"refresh_token_expiry"`
    BCryptCost         int    `mapstructure:"bcrypt_cost"`
}

type RedisConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Password string `mapstructure:"password"`
    DB       int    `mapstructure:"db"`
}

type LoggingConfig struct {
    Level  string `mapstructure:"level"`
    Format string `mapstructure:"format"` // json or console
}
```

### `internal/models/link.go`
```go
package models

import (
    "time"
    
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type Link struct {
    ID           uuid.UUID      `gorm:"type:uuid;primary_key"`
    ShortCode    string         `gorm:"type:varchar(32);uniqueIndex;not null"`
    OriginalURL  string         `gorm:"type:text;not null"`
    UserID       uuid.UUID      `gorm:"type:uuid;index"`
    DomainID     *uuid.UUID     `gorm:"type:uuid;index"`
    Title        string         `gorm:"type:varchar(255)"`
    Description  string         `gorm:"type:text"`
    Tags         pq.StringArray `gorm:"type:text[]"`
    IsActive     bool           `gorm:"default:true"`
    ExpiresAt    *time.Time
    ClickCount   int64          `gorm:"default:0"`
    CreatedAt    time.Time
    UpdatedAt    time.Time
    LastClickedAt *time.Time
    
    // Associations
    User   User    `gorm:"foreignKey:UserID"`
    Domain *Domain `gorm:"foreignKey:DomainID"`
}

func (l *Link) BeforeCreate(tx *gorm.DB) error {
    if l.ID == uuid.Nil {
        l.ID = uuid.New()
    }
    return nil
}

func (l *Link) IsExpired() bool {
    if l.ExpiresAt == nil {
        return false
    }
    return l.ExpiresAt.Before(time.Now())
}
```

### `internal/services/link_service.go`
```go
package services

import (
    "context"
    "errors"
    "time"
    
    "github.com/your-org/url-shortener/internal/models"
    "github.com/your-org/url-shortener/internal/repositories"
    "github.com/your-org/url-shortener/internal/utils"
)

type LinkService interface {
    CreateLink(ctx context.Context, userID uuid.UUID, req CreateLinkRequest) (*models.Link, error)
    GetLink(ctx context.Context, linkID uuid.UUID) (*models.Link, error)
    GetLinkByCode(ctx context.Context, code string) (*models.Link, error)
    ListLinks(ctx context.Context, userID uuid.UUID, filter LinkFilter) ([]models.Link, int64, error)
    UpdateLink(ctx context.Context, linkID uuid.UUID, userID uuid.UUID, updates UpdateLinkRequest) (*models.Link, error)
    DeleteLink(ctx context.Context, linkID uuid.UUID, userID uuid.UUID) error
    IncrementClickCount(ctx context.Context, linkID uuid.UUID) error
}

type linkService struct {
    linkRepo repositories.LinkRepository
    userRepo repositories.UserRepository
    domainRepo repositories.DomainRepository
}

func NewLinkService(
    linkRepo repositories.LinkRepository,
    userRepo repositories.UserRepository,
    domainRepo repositories.DomainRepository,
) LinkService {
    return &linkService{
        linkRepo: linkRepo,
        userRepo: userRepo,
        domainRepo: domainRepo,
    }
}

func (s *linkService) CreateLink(ctx context.Context, userID uuid.UUID, req CreateLinkRequest) (*models.Link, error) {
    // Validate URL
    if !utils.IsValidURL(req.OriginalURL) {
        return nil, errors.New("invalid URL")
    }
    
    // Generate short code if not provided
    code := req.CustomCode
    if code == "" {
        code = utils.GenerateShortCode()
    }
    
    // Check if code is unique
    exists, err := s.linkRepo.CodeExists(ctx, code)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, errors.New("short code already exists")
    }
    
    // Create link
    link := &models.Link{
        ShortCode:   code,
        OriginalURL: req.OriginalURL,
        UserID:      userID,
        Title:       req.Title,
        Description: req.Description,
        Tags:        req.Tags,
        IsActive:    true,
    }
    
    if req.DomainID != nil {
        // Verify user owns the domain
        domain, err := s.domainRepo.GetByID(ctx, *req.DomainID)
        if err != nil {
            return nil, err
        }
        if domain.UserID != userID {
            return nil, errors.New("domain does not belong to user")
        }
        link.DomainID = req.DomainID
    }
    
    if req.ExpiresAt != nil {
        link.ExpiresAt = req.ExpiresAt
    }
    
    err = s.linkRepo.Create(ctx, link)
    if err != nil {
        return nil, err
    }
    
    return link, nil
}
```

### `internal/handlers/link_handler.go`
```go
package handlers

import (
    "net/http"
    
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    
    "github.com/your-org/url-shortener/internal/services"
    "github.com/your-org/url-shortener/internal/utils"
)

type LinkHandler struct {
    linkService services.LinkService
}

func NewLinkHandler(linkService services.LinkService) *LinkHandler {
    return &LinkHandler{linkService: linkService}
}

func (h *LinkHandler) CreateLink(c *gin.Context) {
    var req services.CreateLinkRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "invalid_request", err.Error())
        return
    }
    
    userID, _ := c.Get("user_id")
    uuidUserID, _ := uuid.Parse(userID.(string))
    
    link, err := h.linkService.CreateLink(c.Request.Context(), uuidUserID, req)
    if err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "create_failed", err.Error())
        return
    }
    
    c.JSON(http.StatusCreated, gin.H{
        "data": link,
    })
}

func (h *LinkHandler) GetLink(c *gin.Context) {
    linkID, err := uuid.Parse(c.Param("id"))
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "invalid_id", "Invalid link ID")
        return
    }
    
    userID, _ := c.Get("user_id")
    uuidUserID, _ := uuid.Parse(userID.(string))
    
    link, err := h.linkService.GetLink(c.Request.Context(), linkID)
    if err != nil {
        utils.ErrorResponse(c, http.StatusNotFound, "not_found", "Link not found")
        return
    }
    
    // Authorization check
    if link.UserID != uuidUserID && !c.GetBool("is_admin") {
        utils.ErrorResponse(c, http.StatusForbidden, "forbidden", "You don't have permission to access this link")
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "data": link,
    })
}

func (h *LinkHandler) ListLinks(c *gin.Context) {
    userID, _ := c.Get("user_id")
    uuidUserID, _ := uuid.Parse(userID.(string))
    
    filter := services.LinkFilter{
        Page:     utils.GetQueryInt(c, "page", 1),
        Limit:    utils.GetQueryInt(c, "limit", 20),
        Search:   c.Query("search"),
        Tags:     utils.ParseCommaSeparated(c.Query("tags")),
        ActiveOnly: utils.GetQueryBool(c, "active_only", true),
    }
    
    links, total, err := h.linkService.ListLinks(c.Request.Context(), uuidUserID, filter)
    if err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "list_failed", err.Error())
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "data": links,
        "pagination": gin.H{
            "page":  filter.Page,
            "limit": filter.Limit,
            "total": total,
            "total_pages": (total + filter.Limit - 1) / filter.Limit,
        },
    })
}
```

### `internal/handlers/middleware/auth.go`
```go
package middleware

import (
    "strings"
    
    "github.com/gin-gonic/gin"
    
    "github.com/your-org/url-shortener/internal/utils"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            // Check for API key
            apiKey := c.GetHeader("X-API-Key")
            if apiKey != "" {
                handleAPIKeyAuth(c, apiKey)
                return
            }
            
            c.AbortWithStatusJSON(401, gin.H{
                "error": "authentication_required",
                "message": "Authentication required",
            })
            return
        }
        
        // Bearer token
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.AbortWithStatusJSON(401, gin.H{
                "error": "invalid_token_format",
                "message": "Invalid token format",
            })
            return
        }
        
        token := parts[1]
        claims, err := utils.ValidateJWT(token, jwtSecret)
        if err != nil {
            c.AbortWithStatusJSON(401, gin.H{
                "error": "invalid_token",
                "message": "Invalid or expired token",
            })
            return
        }
        
        // Set user context
        c.Set("user_id", claims.UserID)
        c.Set("user_email", claims.Email)
        c.Set("is_admin", claims.IsAdmin)
        
        c.Next()
    }
}

func handleAPIKeyAuth(c *gin.Context, apiKey string) {
    // Validate API key from database
    // Set user context if valid
    c.Next()
}
```

### `Makefile`
```makefile
.PHONY: help build run test clean migrate seed lint

help:
	@echo "Available commands:"
	@echo "  build     - Build the application"
	@echo "  run       - Run the application"
	@echo "  test      - Run tests"
	@echo "  lint      - Run linter"
	@echo "  migrate   - Run database migrations"
	@echo "  seed      - Seed database with test data"
	@echo "  clean     - Clean build artifacts"

build:
	go build -o bin/url-shortener ./cmd/server

run:
	go run ./cmd/server

test:
	go test ./... -v

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

lint:
	golangci-lint run

migrate-up:
	migrate -path migrations -database "postgres://user:pass@localhost:5432/url_shortener?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://user:pass@localhost:5432/url_shortener?sslmode=disable" down

seed:
	go run scripts/seed.go

clean:
	rm -rf bin/ coverage.out

docker-build:
	docker build -t url-shortener:latest .

docker-run:
	docker-compose up

docker-test:
	docker-compose -f docker-compose.test.yml up --abort-on-container-exit
```

### `docker-compose.yml`
```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: urlshortener
      POSTGRES_PASSWORD: password
      POSTGRES_DB: url_shortener
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./docker/postgres/init.sql:/docker-entrypoint-initdb.d/init.sql
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U urlshortener"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    command: redis-server --appendonly yes

  app:
    build:
      context: .
      dockerfile: docker/Dockerfile.dev
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=urlshortener
      - DB_PASSWORD=password
      - DB_NAME=url_shortener
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - JWT_SECRET=your-secret-key-here
    volumes:
      - .:/app
      - go-mod-cache:/go/pkg/mod
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_started
    command: air

volumes:
  postgres_data:
  redis_data:
  go-mod-cache:
```

### `.env.example`
```env
# Server Configuration
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
SERVER_ENV=development
SERVER_READ_TIMEOUT=30
SERVER_WRITE_TIMEOUT=30
SERVER_IDLE_TIMEOUT=120

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=urlshortener
DB_PASSWORD=password
DB_NAME=url_shortener
DB_SSL_MODE=disable
DB_MAX_CONNS=25

# Authentication
JWT_SECRET=your-super-secret-jwt-key-change-in-production
ACCESS_TOKEN_EXPIRY=900          # 15 minutes in seconds
REFRESH_TOKEN_EXPIRY=604800      # 7 days in seconds
BCRYPT_COST=12

# Redis (Optional)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Logging
LOG_LEVEL=info
LOG_FORMAT=console

# Shortener Configuration
SHORT_CODE_LENGTH=6
BASE_URL=https://short.example.com
ALLOW_CUSTOM_CODES=true
MAX_CUSTOM_CODE_LENGTH=32

# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS=1000
RATE_LIMIT_WINDOW=3600

# Analytics
ANALYTICS_ENABLED=true
ANALYTICS_RETENTION_DAYS=90
ANONYMIZE_IP=true
```

## Development Workflow

### Initial Setup
1. Clone repository
2. Copy `.env.example` to `.env.local` and configure
3. Run `docker-compose up -d postgres redis`
4. Run `make migrate-up`
5. Run `make seed` (optional)
6. Run `make run` or use `docker-compose up app`

### Adding New Features
1. Create model in `internal/models/`
2. Create repository in `internal/repositories/`
3. Create service in `internal/services/`
4. Create handler in `internal/handlers/`
5. Add routes in `internal/app/server.go`
6. Write tests in `tests/`

### Code Organization Principles

1. **Separation of Concerns**: Each layer has distinct responsibilities
2. **Dependency Injection**: Services receive dependencies via interfaces
3. **Error Handling**: Consistent error handling across layers
4. **Testing**: Unit tests for business logic, integration tests for APIs
5. **Documentation**: Code comments for public APIs and complex logic

### Import Aliases
Use consistent import aliases:
- `gorm` for `gorm.io/gorm`
- `gin` for `github.com/gin-gonic/gin`
- `uuid` for `github.com/google/uuid`
- `zap` for `go.uber.org/zap`

### Naming Conventions
- **Files**: `snake_case.go` for implementation, `snake_case_test.go` for tests
- **Packages**: Singular nouns (`user`, `link`, `handler`)
- **Variables**: `camelCase` for locals, `PascalCase` for exports
- **Interfaces**: `Service` suffix for services, `Repository` for repositories

This structure provides a solid foundation for a scalable, maintainable URL shortener service that can grow from MVP to production-ready application.