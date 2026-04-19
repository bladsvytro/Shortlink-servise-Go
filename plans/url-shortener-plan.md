# Go URL Shortener Service - Development Plan

## 1. Project Overview
A full-featured URL shortener service built in Go with PostgreSQL, supporting user authentication, custom domains, API keys, and advanced analytics.

### Core Requirements
- **URL Shortening**: Create short aliases for long URLs
- **Redirection**: Redirect short URLs to original destinations
- **User Management**: Registration, authentication, and authorization
- **Custom Domains**: Allow users to use their own domains
- **API Keys**: Programmatic access with rate limiting
- **Analytics**: Click tracking, geographic data, referral sources
- **Admin Dashboard**: Management interface for administrators

## 2. System Architecture

### High-Level Architecture
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Client Apps   │    │   API Gateway   │    │   Load Balancer │
│   (Web/Mobile)  │───▶│   (Optional)    │───▶│   (Nginx/Traefik)│
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                         │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Redis Cache   │◀───│   Go Service    │───▶│   PostgreSQL    │
│   (Optional)    │    │   (Gin/GORM)    │    │   Database      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                 │
                         ┌─────────────────┐
                         │   Monitoring    │
                         │   (Prometheus)  │
                         └─────────────────┘
```

### Key Components
1. **API Layer**: RESTful endpoints using Gin framework
2. **Business Logic**: URL generation, validation, analytics
3. **Data Layer**: PostgreSQL with GORM ORM
4. **Cache Layer**: Optional Redis for hot paths
5. **Authentication**: JWT-based with refresh tokens
6. **Rate Limiting**: Per-user and per-API-key limits
7. **Analytics Pipeline**: Async processing of click events

## 3. Data Model

### Core Entities

#### User
```go
type User struct {
    ID           uuid.UUID  `gorm:"type:uuid;primary_key"`
    Email        string     `gorm:"uniqueIndex;not null"`
    PasswordHash string     `gorm:"not null"`
    Name         string
    IsActive     bool       `gorm:"default:true"`
    IsAdmin      bool       `gorm:"default:false"`
    CreatedAt    time.Time
    UpdatedAt    time.Time
    LastLoginAt  *time.Time
    APIKeys      []APIKey   `gorm:"foreignKey:UserID"`
    Domains      []Domain   `gorm:"foreignKey:UserID"`
    Links        []Link     `gorm:"foreignKey:UserID"`
}
```

#### Link
```go
type Link struct {
    ID           uuid.UUID  `gorm:"type:uuid;primary_key"`
    ShortCode    string     `gorm:"uniqueIndex;not null"`
    OriginalURL  string     `gorm:"not null"`
    UserID       uuid.UUID  `gorm:"type:uuid;index"`
    DomainID     *uuid.UUID `gorm:"type:uuid;index"`
    Title        string
    Description  string
    Tags         pq.StringArray `gorm:"type:text[]"`
    IsActive     bool           `gorm:"default:true"`
    ExpiresAt    *time.Time
    ClickCount   int64          `gorm:"default:0"`
    CreatedAt    time.Time
    UpdatedAt    time.Time
    LastClickedAt *time.Time
}
```

#### Domain
```go
type Domain struct {
    ID           uuid.UUID  `gorm:"type:uuid;primary_key"`
    DomainName   string     `gorm:"uniqueIndex;not null"`
    UserID       uuid.UUID  `gorm:"type:uuid;index"`
    IsVerified   bool       `gorm:"default:false"`
    IsActive     bool       `gorm:"default:true"`
    CreatedAt    time.Time
    VerifiedAt   *time.Time
}
```

#### APIKey
```go
type APIKey struct {
    ID           uuid.UUID  `gorm:"type:uuid;primary_key"`
    UserID       uuid.UUID  `gorm:"type:uuid;index"`
    KeyHash      string     `gorm:"uniqueIndex;not null"`
    Name         string
    LastUsedAt   *time.Time
    ExpiresAt    *time.Time
    RateLimit    int        `gorm:"default:1000"`
    CreatedAt    time.Time
}
```

#### ClickEvent
```go
type ClickEvent struct {
    ID           uuid.UUID  `gorm:"type:uuid;primary_key"`
    LinkID       uuid.UUID  `gorm:"type:uuid;index"`
    IPAddress    string
    UserAgent    string
    Referrer     string
    CountryCode  string
    City         string
    DeviceType   string     // mobile, desktop, tablet
    Browser      string
    OS           string
    Timestamp    time.Time
}
```

### Database Schema Relationships
- One User → Many Links
- One User → Many Domains  
- One User → Many APIKeys
- One Domain → Many Links
- One Link → Many ClickEvents

## 4. API Endpoints

### Authentication Endpoints
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login (JWT)
- `POST /api/v1/auth/refresh` - Refresh JWT token
- `POST /api/v1/auth/logout` - Invalidate token
- `GET /api/v1/auth/profile` - Get user profile

### Link Management Endpoints
- `POST /api/v1/links` - Create new short link
- `GET /api/v1/links` - List user's links (with pagination)
- `GET /api/v1/links/{id}` - Get link details
- `PUT /api/v1/links/{id}` - Update link
- `DELETE /api/v1/links/{id}` - Delete link
- `GET /api/v1/links/{id}/analytics` - Get link analytics
- `GET /api/v1/links/{id}/clicks` - Get click events (with filters)

### Domain Management Endpoints
- `POST /api/v1/domains` - Add custom domain
- `GET /api/v1/domains` - List user's domains
- `DELETE /api/v1/domains/{id}` - Remove domain
- `POST /api/v1/domains/{id}/verify` - Verify domain ownership

### API Key Management
- `POST /api/v1/apikeys` - Generate new API key
- `GET /api/v1/apikeys` - List API keys
- `DELETE /api/v1/apikeys/{id}` - Revoke API key

### Redirect Endpoint
- `GET /{shortCode}` - Redirect to original URL (public)
- `GET /{domain}/{shortCode}` - Redirect with custom domain

### Admin Endpoints
- `GET /admin/users` - List all users
- `PUT /admin/users/{id}/status` - Toggle user status
- `GET /admin/links` - List all links
- `GET /admin/stats` - System statistics

## 5. Technology Stack

### Core Framework & Libraries
- **Web Framework**: Gin (lightweight, high performance)
- **ORM**: GORM (database abstraction)
- **Validation**: go-playground/validator
- **Configuration**: Viper (config management)
- **Logging**: Zap (structured logging)
- **JWT**: golang-jwt/jwt
- **UUID**: google/uuid
- **Migration**: golang-migrate/migrate

### Database
- **Primary**: PostgreSQL 14+
- **Connection Pool**: pgx
- **Migration Tool**: golang-migrate

### Optional Components
- **Cache**: Redis (for rate limiting and hot paths)
- **Message Queue**: RabbitMQ/Kafka (for async analytics)
- **Search**: Elasticsearch (for advanced analytics queries)
- **Object Storage**: MinIO/S3 (for file uploads if needed)

### Monitoring & Observability
- **Metrics**: Prometheus + Grafana
- **Tracing**: OpenTelemetry
- **Health Checks**: Health endpoint with dependencies

## 6. Project Structure

```
url-shortener/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── config/              # Configuration management
│   │   └── config.go
│   ├── database/            # Database connection and migrations
│   │   ├── connection.go
│   │   └── migrations/      # SQL migration files
│   ├── models/              # Data models (GORM structs)
│   │   ├── user.go
│   │   ├── link.go
│   │   └── ...
│   ├── repositories/        # Data access layer
│   │   ├── user_repo.go
│   │   ├── link_repo.go
│   │   └── ...
│   ├── services/            # Business logic
│   │   ├── auth_service.go
│   │   ├── link_service.go
│   │   ├── analytics_service.go
│   │   └── ...
│   ├── handlers/            # HTTP handlers (controllers)
│   │   ├── auth_handler.go
│   │   ├── link_handler.go
│   │   ├── domain_handler.go
│   │   └── ...
│   ├── middleware/          # HTTP middleware
│   │   ├── auth.go
│   │   ├── logging.go
│   │   ├── rate_limit.go
│   │   └── ...
│   ├── utils/               # Utility functions
│   │   ├── jwt.go
│   │   ├── validator.go
│   │   └── ...
│   └── pkg/                 # Reusable packages (if any)
├── api/
│   └── v1/                  # API specifications
│       ├── openapi.yaml     # OpenAPI specification
│       └── client/          # Generated API client
├── migrations/              # Database migration files
│   ├── 001_init_schema.up.sql
│   └── 001_init_schema.down.sql
├── tests/                   # Test files
│   ├── unit/
│   └── integration/
├── docker/                  # Docker configurations
│   ├── Dockerfile
│   └── docker-compose.yml
├── deployments/             # Deployment configurations
│   ├── kubernetes/
│   └── helm/
├── scripts/                 # Utility scripts
│   ├── migrate.sh
│   └── seed.sh
├── .env.example             # Environment variables template
├── .gitignore
├── go.mod                   # Go module definition
├── go.sum
├── Makefile                 # Common tasks
└── README.md
```

## 7. Database Schema Details

### Tables
1. **users** - User accounts and profiles
2. **links** - Shortened URLs and metadata
3. **domains** - Custom domains for URL shortening
4. **api_keys** - API keys for programmatic access
5. **click_events** - Analytics data for each click
6. **sessions** - User sessions (optional, for JWT blacklist)

### Indexes
- `links.short_code` (unique) - For fast redirect lookups
- `links.user_id` - For user-specific queries
- `click_events.link_id` + `timestamp` - For time-based analytics
- `users.email` (unique) - For authentication

## 8. Authentication & Authorization

### Authentication Flow
1. User registers with email/password
2. System hashes password (bcrypt) and creates user
3. User logs in with credentials
4. System validates and issues JWT access token (15min) + refresh token (7 days)
5. Access token used for API calls, refresh token for obtaining new access tokens

### Authorization Levels
1. **Public**: Redirect endpoints only
2. **User**: Own resources (links, domains, API keys)
3. **Admin**: All resources, system management

### API Key Authentication
- API keys passed in `X-API-Key` header
- Rate limiting per key (configurable)
- Keys can be revoked anytime

## 9. Key Algorithms & Logic

### Short Code Generation
```go
func GenerateShortCode() string {
    // Base62 encoding of random bytes or timestamp
    // Options: random, sequential, or custom
    // Ensure uniqueness via database check
}
```

### URL Validation
- Check URL format and scheme
- Ensure URL is reachable (optional)
- Prevent malicious URLs (phishing detection)

### Rate Limiting
- Token bucket algorithm
- Per user and per API key limits
- Redis-backed for distributed systems

## 10. Monitoring & Analytics

### Metrics to Track
- Total links created
- Total redirects (per link, per user, system-wide)
- API request rates and errors
- Database connection pool status
- Response time percentiles

### Analytics Features
1. **Basic**: Click counts over time
2. **Geographic**: Map of clicks by country/city
3. **Referral**: Traffic sources
4. **Device/Browser**: User agent analysis
5. **Time Series**: Hourly/daily/weekly trends

## 11. Deployment Considerations

### Development
- Docker Compose for local development
- Hot reload with Air/CompileDaemon
- Seed data for testing

### Production
- Containerized deployment (Docker)
- Kubernetes for orchestration
- PostgreSQL with replication
- Load balancing with Nginx/Traefik
- CI/CD pipeline (GitHub Actions/GitLab CI)

### Environment Configuration
- 12-factor app principles
- Environment variables for secrets
- Config files for non-secret settings

## 12. Security Considerations

### Input Validation
- Validate all user inputs
- Sanitize URLs to prevent XSS
- Rate limit all endpoints

### Data Protection
- HTTPS mandatory
- Password hashing with bcrypt
- JWT signing with strong secret
- Database encryption at rest

### Privacy
- Anonymize IP addresses in analytics (optional)
- GDPR compliance considerations
- Data retention policies

## 13. Testing Strategy

### Test Types
1. **Unit Tests**: Business logic, utilities
2. **Integration Tests**: Database operations, API endpoints
3. **E2E Tests**: Full user flows
4. **Load Tests**: Performance under stress

### Test Tools
- Go's built-in testing package
- Testify for assertions
- Mockery for mocking
- GoConvey for BDD (optional)

## 14. Development Roadmap

### Phase 1: MVP (Weeks 1-2)
- Basic URL shortening without auth
- Redirect functionality
- Simple analytics (click count)
- PostgreSQL setup

### Phase 2: Core Features (Weeks 3-4)
- User authentication (JWT)
- Link management API
- Basic dashboard
- API key support

### Phase 3: Advanced Features (Weeks 5-6)
- Custom domains
- Advanced analytics
- Rate limiting
- Admin panel

### Phase 4: Polish & Scale (Weeks 7-8)
- Caching (Redis)
- Async analytics processing
- Monitoring & alerts
- Documentation

## 15. Success Metrics

### Technical Metrics
- API response time < 100ms (p95)
- Redirect time < 50ms (p95)
- 99.9% availability
- Zero data loss

### Business Metrics
- User adoption rate
- Links created per user
- Click-through rates
- Domain verification rate

---

*This plan will be refined based on feedback and additional requirements.*