# Monitoring, Logging, and Deployment Plan

## Observability Stack

### 1. Logging Strategy

#### Structured Logging with Zap
```go
// Logger configuration
func NewLogger(level string, format string) (*zap.Logger, error) {
    config := zap.NewProductionConfig()
    
    // Set log level
    switch level {
    case "debug":
        config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
    case "info":
        config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
    case "warn":
        config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
    case "error":
        config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
    }
    
    // Set output format
    if format == "console" {
        config.Encoding = "console"
        config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
    } else {
        config.Encoding = "json"
    }
    
    // Add custom fields
    config.InitialFields = map[string]interface{}{
        "service": "url-shortener",
        "version": "1.0.0",
    }
    
    return config.Build()
}

// Logging middleware
func LoggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        // Process request
        c.Next()
        
        // Log after request
        latency := time.Since(start)
        
        logger.Info("HTTP request",
            zap.String("method", c.Request.Method),
            zap.String("path", c.Request.URL.Path),
            zap.Int("status", c.Writer.Status()),
            zap.Duration("latency", latency),
            zap.String("client_ip", c.ClientIP()),
            zap.String("user_agent", c.Request.UserAgent()),
            zap.String("request_id", c.GetString("request_id")),
            zap.String("user_id", c.GetString("user_id")),
        )
    }
}
```

#### Log Levels and Usage
| Level | When to Use | Example |
|-------|-------------|---------|
| DEBUG | Development, troubleshooting | `logger.Debug("Processing link", zap.String("code", code))` |
| INFO | Normal operations, request logging | `logger.Info("Link created", zap.String("code", code))` |
| WARN | Unexpected but recoverable situations | `logger.Warn("Rate limit approaching", zap.String("user", userID))` |
| ERROR | Failures that need attention | `logger.Error("Database connection failed", zap.Error(err))` |
| FATAL | Critical failures requiring shutdown | `logger.Fatal("Failed to bind port", zap.Error(err))` |

#### Log Aggregation
- **Development**: Console output with colors
- **Production**: JSON format to stdout (container logs)
- **Aggregation**: Fluentd/Fluent Bit → Elasticsearch/Loki
- **Retention**: 30 days for debug/info, 90 days for error/fatal

### 2. Metrics and Monitoring

#### Prometheus Metrics
```go
// Metrics setup
func SetupMetrics() {
    // Custom metrics
    linksCreated = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "url_shortener_links_created_total",
            Help: "Total number of links created",
        },
        []string{"user_id"},
    )
    
    redirectsServed = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "url_shortener_redirects_served_total",
            Help: "Total number of redirects served",
        },
        []string{"code", "status"},
    )
    
    apiRequests = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "url_shortener_api_requests_total",
            Help: "Total API requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "url_shortener_request_duration_seconds",
            Help:    "Request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "endpoint"},
    )
    
    // Register metrics
    prometheus.MustRegister(linksCreated, redirectsServed, apiRequests, requestDuration)
}

// Metrics middleware
func MetricsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        c.Next()
        
        duration := time.Since(start).Seconds()
        status := strconv.Itoa(c.Writer.Status())
        
        apiRequests.WithLabelValues(
            c.Request.Method,
            c.FullPath(),
            status,
        ).Inc()
        
        requestDuration.WithLabelValues(
            c.Request.Method,
            c.FullPath(),
        ).Observe(duration)
    }
}
```

#### Key Metrics to Monitor
```yaml
# Business Metrics
- links_created_total
- redirects_served_total
- active_users_total
- domains_verified_total

# Performance Metrics
- request_duration_seconds
- api_requests_total
- database_query_duration_seconds
- cache_hit_ratio

# System Metrics
- go_goroutines
- go_memstats_alloc_bytes
- process_cpu_seconds_total
- process_resident_memory_bytes

# Database Metrics
- postgres_connections_total
- postgres_queries_per_second
- postgres_replication_lag

# External Dependencies
- redis_connected
- redis_command_duration_seconds
```

#### Grafana Dashboards
1. **Service Overview Dashboard**
   - Request rate, error rate, latency (p50, p95, p99)
   - CPU, memory, goroutine count
   - Database connection pool status

2. **Business Metrics Dashboard**
   - Links created per hour/day
   - Redirects served (top links, geographic distribution)
   - User growth, active users

3. **Database Dashboard**
   - Query performance, slow queries
   - Connection pool utilization
   - Replication lag (if applicable)

### 3. Tracing and Distributed Context

#### OpenTelemetry Setup
```go
func SetupTracing(serviceName string) (*sdktrace.TracerProvider, error) {
    // Create exporter (Jaeger, Zipkin, or OTLP)
    exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(
        jaeger.WithEndpoint("http://jaeger:14268/api/traces"),
    ))
    if err != nil {
        return nil, err
    }
    
    tp := sdktrace.NewTracerProvider(
        sdktrace.WithSampler(sdktrace.AlwaysSample()),
        sdktrace.WithBatcher(exporter),
        sdktrace.WithResource(resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceNameKey.String(serviceName),
            semconv.ServiceVersionKey.String("1.0.0"),
        )),
    )
    
    otel.SetTracerProvider(tp)
    return tp, nil
}

// Tracing middleware
func TracingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tracer := otel.Tracer("url-shortener")
        ctx, span := tracer.Start(c.Request.Context(), c.FullPath())
        defer span.End()
        
        // Add trace context to request
        c.Request = c.Request.WithContext(ctx)
        
        // Set trace ID in response headers
        spanContext := span.SpanContext()
        if spanContext.HasTraceID() {
            c.Header("X-Trace-ID", spanContext.TraceID().String())
        }
        
        c.Next()
    }
}
```

### 4. Health Checks

#### Health Check Endpoints
```go
func SetupHealthChecks(db *gorm.DB, redisClient *redis.Client) {
    health := health.New()
    
    // Database health check
    health.AddReadinessCheck("database", healthcheck.DatabasePingCheck(db.DB(), 2*time.Second))
    
    // Redis health check (if used)
    if redisClient != nil {
        health.AddReadinessCheck("redis", func() error {
            ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
            defer cancel()
            return redisClient.Ping(ctx).Err()
        })
    }
    
    // Custom health checks
    health.AddLivenessCheck("goroutine-threshold", goroutineCountCheck(1000))
    
    // HTTP handlers
    http.HandleFunc("/health/live", health.LiveEndpoint)
    http.HandleFunc("/health/ready", health.ReadyEndpoint)
}

func goroutineCountCheck(threshold int) healthcheck.Check {
    return func() error {
        count := runtime.NumGoroutine()
        if count > threshold {
            return fmt.Errorf("too many goroutines: %d", count)
        }
        return nil
    }
}
```

#### Kubernetes Probes
```yaml
# Deployment configuration
livenessProbe:
  httpGet:
    path: /health/live
    port: 8080
  initialDelaySeconds: 30
  periodSeconds: 10
  timeoutSeconds: 5
  failureThreshold: 3

readinessProbe:
  httpGet:
    path: /health/ready
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
  timeoutSeconds: 3
  failureThreshold: 1
```

## Deployment Strategy

### 1. Containerization

#### Dockerfile (Production)
```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o url-shortener ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/url-shortener .
COPY --from=builder /app/migrations ./migrations

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

USER 1001

# Expose port
EXPOSE 8080

# Run application
CMD ["./url-shortener"]
```

#### Dockerfile (Development)
```dockerfile
FROM golang:1.21-alpine

WORKDIR /app

# Install air for hot reload
RUN go install github.com/cosmtrek/air@latest

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Expose port
EXPOSE 8080

# Run with air for hot reload
CMD ["air", "-c", ".air.toml"]
```

### 2. Docker Compose Setup

#### docker-compose.yml
```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: ${DB_USER:-urlshortener}
      POSTGRES_PASSWORD: ${DB_PASSWORD:-password}
      POSTGRES_DB: ${DB_NAME:-url_shortener}
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./docker/postgres/init.sql:/docker-entrypoint-initdb.d/init.sql
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER:-urlshortener}"]
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
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  app:
    build:
      context: .
      dockerfile: Dockerfile.dev
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=${DB_USER:-urlshortener}
      - DB_PASSWORD=${DB_PASSWORD:-password}
      - DB_NAME=${DB_NAME:-url_shortener}
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - JWT_SECRET=${JWT_SECRET:-your-secret-key}
      - LOG_LEVEL=debug
    volumes:
      - .:/app
      - go-mod-cache:/go/pkg/mod
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    command: air

  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./docker/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--web.console.libraries=/etc/prometheus/console_libraries'
      - '--web.console.templates=/etc/prometheus/consoles'
      - '--storage.tsdb.retention.time=200h'
      - '--web.enable-lifecycle'

  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=${GRAFANA_PASSWORD:-admin}
    volumes:
      - grafana_data:/var/lib/grafana
      - ./docker/grafana/dashboards:/etc/grafana/provisioning/dashboards
      - ./docker/grafana/datasources:/etc/grafana/provisioning/datasources
    depends_on:
      - prometheus

volumes:
  postgres_data:
  redis_data:
  prometheus_data:
  grafana_data:
  go-mod-cache:
```

### 3. Kubernetes Deployment

#### Deployment Configuration
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: url-shortener
  namespace: production
spec:
  replicas: 3
  selector:
    matchLabels:
      app: url-shortener
  template:
    metadata:
      labels:
        app: url-shortener
    spec:
      containers:
      - name: app
        image: your-registry/url-shortener:latest
        ports:
        - containerPort: 8080
        env:
        - name: DB_HOST
          valueFrom:
            secretKeyRef:
              name: database-secrets
              key: host
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: database-secrets
              key: password
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: jwt-secret
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "200m"
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

#### Service Configuration
```yaml
apiVersion: v1
kind: Service
metadata:
  name: url-shortener
  namespace: production
spec:
  selector:
    app: url-shortener
  ports:
  - port: 80
    targetPort: 8080
  type: ClusterIP
```

#### Ingress Configuration
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: url-shortener
  namespace: production
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
spec:
  tls:
  - hosts:
    - api.short.example.com
    - short.example.com
    secretName: url-shortener-tls
  rules:
  - host: api.short.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: url-shortener
            port:
              number: 80
  - host: short.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: url-shortener
            port:
              number: 80
```

### 4. CI/CD Pipeline

#### GitHub Actions Workflow
```yaml
name: CI/CD Pipeline

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    
    - name: Run tests
      run: go test ./... -v -coverprofile=coverage.out
    
    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.out

  lint:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Run golangci-lint
      uses: golangci/golangci-lint-action@v3
      with:
        version: latest

  build:
    needs: [test, lint]
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Docker Buildx
      uses: docker/setup-buildx-action@v2
    
    - name: Login to DockerHub
      uses: docker/login-action@v2
      with:
        username: ${{ secrets.DOCKER_USERNAME }}
        password: ${{ secrets.DOCKER_PASSWORD }}
    
    - name: Build and push
      uses: docker/build-push-action@v4
      with:
        context: .
        push: true
        tags: |
          your-registry/url-shortener:latest
          your-registry/url-shortener:${{ github.sha }}

  deploy:
    needs: build
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    steps:
    - name: Deploy to Kubernetes
      uses: azure/k8s-deploy@v4
      with:
        namespace: production
        manifests: |
          deployments/kubernetes/deployment.yaml
          deployments/kubernetes/service.yaml
          deployments/kubernetes/ingress.yaml
        images: |
          your-registry/url-shortener:${{ github.sha }}
```

### 5. Environment Configuration

#### Environment Variables
```bash
# Required
DB_HOST=localhost
DB_PORT=5432
DB_USER=urlshortener
DB_PASSWORD=password
DB_NAME=url_shortener
JWT_SECRET=your-super-secret-jwt-key

# Optional with defaults
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
LOG_LEVEL=info
LOG_FORMAT=json
REDIS_HOST=localhost
REDIS_PORT=6379
BASE_URL=https://short.example.com
```

#### Configuration Management
```go
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Auth     AuthConfig
    Redis    RedisConfig
    Logging  LoggingConfig
    Shortener ShortenerConfig
}

func Load() (*Config, error) {
    // Load from .env file
    _ = godotenv.Load(".env.local")
    
    var config Config
    
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    viper.AddConfigPath("./config")
    
    // Read config file
    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            return nil, err
        }
    }
    
    // Bind environment variables
    viper.AutomaticEnv()
    viper.SetEnvPrefix("URL_SHORTENER")
    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
    
    // Unmarshal config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, err
    }
    
    return &config, nil
}
```

### 6. Backup and Disaster Recovery

#### Database Backups
```bash
#!/bin/bash
# backup.sh

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backups"
FILENAME="url_shortener_${DATE}.dump"

# Create backup
pg_dump -U $DB_USER -h $DB_HOST -d $DB_NAME \
  --format=custom --blobs --verbose \
  --file="${BACKUP_DIR}/${FILENAME}"

# Upload to S3 (optional)
aws s3 cp "${BACKUP_DIR}/${FILENAME}" "s3://your-backup-bucket/${FILENAME}"

# Cleanup old backups (keep 30 days)
find $BACKUP_DIR -name "*.dump" -mtime +30 -delete
```

#### Recovery Procedure
1. **Database corruption**: Restore from latest backup
2. **Data loss**: Point-in-time recovery using WAL archives
3. **Service outage**: Failover to standby database
4. **Configuration error**: Rollback to previous deployment

### 7. Scaling Strategy

#### Horizontal Scaling
```yaml
# Horizontal Pod Autoscaler
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: url-shortener
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: url-shortener
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

#### Database Scaling
1. **Read replicas** for analytics queries
2. **Connection pooling** with PgBouncer
3. **Query optimization** and indexing
4. **Partitioning** for click_events table

### 8. Security Considerations

#### Network Security
```yaml
# NetworkPolicy
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: url-shortener-policy
spec:
  podSelector:
    matchLabels:
      app: url-shortener
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          name: ingress-nginx
    ports:
    - protocol: TCP
      port: 8080
  egress:
  - to:
    - podSelector:
        matchLabels:
          app: postgres
    ports:
    - protocol: TCP
      port: 5432
```

#### Secret Management
```bash
# Create Kubernetes secrets
kubectl create secret generic database-secrets \
  --from-literal=host=postgres \
  --from-literal=password=$(openssl rand -base64 32)

kubectl create secret generic app-secrets \
  --from-literal=jwt-secret=$(openssl rand -base64 64)
```

### 9. Monitoring and Alerting

#### Alert Rules (Prometheus)
```yaml
groups:
- name: url-shortener
  rules:
  - alert: HighErrorRate
    expr: rate(url_shortener_api_requests_total{status=~"5.."}[5m]) / rate(url_shortener_api_requests_total[5m]) > 0.05
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: "High error rate detected"
      description: "Error rate is above 5% for 5 minutes"
  
  - alert: HighLatency
    expr: histogram_quantile(0.95, rate(url_shortener_request_duration_seconds_bucket[5m])) > 1
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "High latency detected"
      description: "95th percentile latency is above 1 second"
  
  - alert: ServiceDown
    expr: up{job="url-shortener"} == 0
    for: 1m
    labels:
      severity: critical
    annotations:
      summary: "Service is down"
      description: "URL Shortener service is not responding"
```

#### Notification Channels
- **Slack**: Critical alerts
- **Email**: Daily/weekly reports
- **PagerDuty**: On-call alerts
- **Webhook**: Custom integrations

This comprehensive monitoring, logging, and deployment plan ensures the URL shortener service is observable, maintainable, and scalable in production environments.