# Technology Stack for Go URL Shortener

## Core Dependencies

### Go Version
- **Go 1.21+** (or latest stable)
- **Module**: `go.mod` with minimum version selection

### Web Framework
- **Gin** (`github.com/gin-gonic/gin`) - High-performance HTTP web framework
  - Lightweight, fast routing
  - Middleware support
  - JSON binding and validation
  - Recovery from panics
- **Alternative considered**: Echo, Fiber, Chi

### Database ORM
- **GORM** (`gorm.io/gorm`) - ORM library for Go
  - Supports PostgreSQL, MySQL, SQLite
  - Migrations, associations, hooks
  - Connection pooling, transactions
- **Driver**: `gorm.io/driver/postgres` - PostgreSQL driver

### Database
- **PostgreSQL 14+** - Primary relational database
  - JSONB support for flexible data
  - Full-text search capabilities
  - Geospatial extensions (optional for analytics)
- **Connection Pool**: `github.com/jackc/pgx/v5` - High-performance PostgreSQL driver

### Authentication & Security
- **JWT**: `github.com/golang-jwt/jwt/v5` - JSON Web Tokens
- **Password Hashing**: `golang.org/x/crypto/bcrypt` - Secure password hashing
- **UUID**: `github.com/google/uuid` - UUID generation
- **Validation**: `github.com/go-playground/validator/v10` - Struct validation

### Configuration
- **Viper** (`github.com/spf13/viper`) - Configuration management
  - Support for JSON, YAML, TOML, env vars
  - Live reloading (optional)
- **Env**: `github.com/joho/godotenv` - Load environment variables from .env

### Logging
- **Zap** (`go.uber.org/zap`) - High-performance structured logging
  - JSON and console output
  - Log levels, structured fields
  - Performance optimized
- **Alternative**: Logrus, Zerolog

### Testing
- **Testing Framework**: Go's built-in `testing` package
- **Assertions**: `github.com/stretchr/testify` - Test assertions and mocks
  - `assert` for assertions
  - `require` for required assertions
  - `mock` for mocking interfaces
- **HTTP Testing**: `net/http/httptest`
- **Test Containers**: `github.com/testcontainers/testcontainers-go` (optional for integration tests)

### Migration Tool
- **golang-migrate** (`github.com/golang-migrate/migrate`) - Database migrations
  - SQL-based migrations (up/down)
  - Version tracking
  - CLI tool included

### HTTP Client
- **Standard Library**: `net/http` (for external API calls)
- **Enhanced**: `github.com/go-resty/resty/v2` (optional for complex HTTP requests)

### Utilities
- **String Manipulation**: Standard library + `strings`, `strconv`
- **Time**: `time` package with `github.com/itchyny/timefmt-go` for formatting
- **CSV/JSON**: Standard encoding packages
- **Compression**: `compress/gzip` for response compression

## Optional Dependencies

### Caching (if needed)
- **Redis**: `github.com/go-redis/redis/v8` - Redis client
- **In-memory**: `github.com/patrickmn/go-cache` - Simple in-memory cache

### Message Queue (for async processing)
- **RabbitMQ**: `github.com/streadway/amqp`
- **NATS**: `github.com/nats-io/nats.go`
- **In-process**: `github.com/hibiken/asynq` (Redis-based)

### Monitoring & Metrics
- **Prometheus**: `github.com/prometheus/client_golang` - Metrics collection
- **Health Checks**: `github.com/heptiolabs/healthcheck`
- **Tracing**: `go.opentelemetry.io/otel` - OpenTelemetry support

### Email Service (for notifications)
- **SMTP**: Standard `net/smtp`
- **Enhanced**: `github.com/go-gomail/gomail` or `github.com/wneessen/go-mail`

### File Storage (if needed for exports)
- **AWS S3**: `github.com/aws/aws-sdk-go-v2/service/s3`
- **Local**: Standard `os` package

### API Documentation
- **Swagger/OpenAPI**: `github.com/swaggo/swag` + `github.com/swaggo/gin-swagger`
- **Alternative**: `github.com/getkin/kin-openapi`

## Development Tools

### Code Quality
- **Linter**: `golangci-lint` - Aggregated linters
- **Formatter**: `gofmt` (built-in) + `goimports`
- **Static Analysis**: `staticcheck`, `govulncheck`

### Hot Reload (Development)
- **Air** (`github.com/cosmtrek/air`) - Live reload for Go apps
- **CompileDaemon** (`github.com/githubnemo/CompileDaemon`) - Alternative

### Code Generation
- **Mock Generation**: `github.com/vektra/mockery/v2` - Generate mocks
- **Wire**: `github.com/google/wire` - Dependency injection (optional)

### Build & Deployment
- **Makefile** - Task automation
- **Docker** - Containerization
- **GitHub Actions/GitLab CI** - CI/CD pipelines

## Project Dependencies (go.mod example)

```go
module github.com/your-org/url-shortener

go 1.21

require (
    // Web Framework
    github.com/gin-gonic/gin v1.9.1
    
    // Database
    gorm.io/gorm v1.25.5
    gorm.io/driver/postgres v1.5.4
    github.com/jackc/pgx/v5 v5.5.0
    
    // Authentication & Security
    github.com/golang-jwt/jwt/v5 v5.2.0
    golang.org/x/crypto v0.17.0
    github.com/google/uuid v1.5.0
    
    // Validation
    github.com/go-playground/validator/v10 v10.16.0
    
    // Configuration
    github.com/spf13/viper v1.18.0
    github.com/joho/godotenv v1.5.1
    
    // Logging
    go.uber.org/zap v1.26.0
    
    // Testing
    github.com/stretchr/testify v1.8.4
    
    // Migrations
    github.com/golang-migrate/migrate/v4 v4.16.2
    
    // Utilities
    github.com/itchyny/timefmt-go v0.1.5
)

// Development dependencies
require (
    github.com/cosmtrek/air v1.49.0 // dev
    github.com/golangci/golangci-lint v1.55.2 // dev
)
```

## Architecture Decisions

### Why Gin over other frameworks?
- **Performance**: Gin is one of the fastest Go web frameworks
- **Middleware**: Rich middleware ecosystem
- **Community**: Large community, well-documented
- **Simplicity**: Easy to learn and use

### Why GORM over raw SQL?
- **Productivity**: Faster development with struct mapping
- **Type Safety**: Compile-time checking of queries
- **Migrations**: Built-in migration support
- **Relationships**: Easy handling of associations

### Why PostgreSQL over other databases?
- **Reliability**: ACID compliance, data integrity
- **JSON Support**: JSONB for flexible schema
- **Full-text Search**: Built-in search capabilities
- **Extensions**: Rich ecosystem of extensions
- **Analytics**: Window functions, CTEs for complex queries

### Why Zap for logging?
- **Performance**: Zero-allocation design
- **Structured Logging**: JSON format for log aggregation
- **Levels**: Fine-grained log levels
- **Integration**: Works well with monitoring systems

### Why JWT for authentication?
- **Stateless**: No server-side session storage
- **Scalability**: Easy to scale horizontally
- **Standard**: Widely adopted standard
- **Flexibility**: Can include custom claims

## Deployment Considerations

### Containerization
- **Base Image**: `golang:1.21-alpine` for build, `alpine:latest` for runtime
- **Multi-stage Builds**: Reduce final image size
- **Non-root User**: Run as non-root for security

### Environment Configuration
- **12-Factor App**: Configuration via environment variables
- **Secrets Management**: External secrets (Kubernetes Secrets, AWS Secrets Manager)
- **Feature Flags**: Configuration for enabling/disabling features

### Health Checks
- **Readiness Probe**: `/health/ready` - checks database connection
- **Liveness Probe**: `/health/live` - basic application health
- **Metrics Endpoint**: `/metrics` - Prometheus metrics

### Scaling
- **Horizontal Scaling**: Stateless design allows multiple instances
- **Database Connection Pool**: Configure appropriate pool size
- **Caching Layer**: Redis for frequently accessed data

## Monitoring Stack

### Required
- **Application Metrics**: Prometheus metrics endpoint
- **Log Aggregation**: Structured logs to centralized system
- **Error Tracking**: Sentry or similar (optional)

### Optional
- **Distributed Tracing**: OpenTelemetry for request tracing
- **Performance Monitoring**: APM tools (Datadog, New Relic)
- **Business Metrics**: Custom metrics for business KPIs

## Security Considerations

### Dependencies
- **Regular Updates**: Use Dependabot/Renovate for dependency updates
- **Vulnerability Scanning**: `govulncheck` and Snyk/Trivy
- **Minimal Dependencies**: Only include necessary packages

### Application Security
- **Input Validation**: Validate all user inputs
- **SQL Injection Prevention**: Use parameterized queries (GORM handles this)
- **XSS Protection**: Sanitize URLs and user content
- **Rate Limiting**: Prevent abuse
- **CORS**: Configure appropriate CORS policies

## Development Workflow

### Local Development
1. `go mod download` - Install dependencies
2. `air` - Start development server with hot reload
3. `docker-compose up` - Start PostgreSQL and other services
4. `make migrate-up` - Run database migrations
5. `make test` - Run tests

### CI/CD Pipeline
1. **Lint**: `golangci-lint run`
2. **Test**: `go test ./...`
3. **Build**: `go build -o url-shortener`
4. **Security Scan**: Vulnerability scanning
5. **Container Build**: Build Docker image
6. **Deploy**: Deploy to staging/production

## Alternative Stack Options

### Simplified Stack (for smaller projects)
- **Framework**: Standard `net/http` or Chi
- **Database**: SQLite (for single-instance deployment)
- **Authentication**: Basic auth or simple API keys
- **Logging**: Standard `log` package

### Enterprise Stack (for large scale)
- **Framework**: Gin with custom middleware
- **Database**: PostgreSQL with read replicas
- **Cache**: Redis cluster
- **Message Queue**: Kafka for event streaming
- **Search**: Elasticsearch for analytics
- **Monitoring**: Full observability stack

## Version Compatibility

All selected libraries are compatible with Go 1.21+ and actively maintained. Regular updates should be scheduled to keep dependencies current and secure.

## Performance Considerations

- **Connection Pooling**: Configure database connection pool appropriately
- **Query Optimization**: Use database indexes, avoid N+1 queries
- **Caching**: Implement caching for frequently accessed data
- **Compression**: Enable gzip compression for API responses
- **Concurrency**: Leverage Go's goroutines for parallel processing where appropriate