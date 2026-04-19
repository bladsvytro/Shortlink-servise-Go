# Database Schema and Migrations

## PostgreSQL Schema Design

### Tables Overview

1. **users** - User accounts and authentication
2. **links** - Shortened URLs and metadata  
3. **domains** - Custom domains for URL shortening
4. **api_keys** - API keys for programmatic access
5. **click_events** - Analytics data for each click
6. **sessions** - User sessions (optional, for JWT blacklist)
7. **migrations** - Migration version tracking (managed by golang-migrate)

## Detailed Table Schemas

### users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    is_active BOOLEAN DEFAULT TRUE,
    is_admin BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_login_at TIMESTAMP WITH TIME ZONE,
    
    -- Indexes
    INDEX idx_users_email (email),
    INDEX idx_users_created_at (created_at),
    INDEX idx_users_is_active (is_active)
);

-- Trigger for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
```

### links Table
```sql
CREATE TABLE links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    short_code VARCHAR(32) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    domain_id UUID REFERENCES domains(id) ON DELETE SET NULL,
    title VARCHAR(255),
    description TEXT,
    tags TEXT[] DEFAULT '{}',
    is_active BOOLEAN DEFAULT TRUE,
    expires_at TIMESTAMP WITH TIME ZONE,
    click_count BIGINT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_clicked_at TIMESTAMP WITH TIME ZONE,
    
    -- Indexes
    INDEX idx_links_short_code (short_code),
    INDEX idx_links_user_id (user_id),
    INDEX idx_links_domain_id (domain_id),
    INDEX idx_links_created_at (created_at),
    INDEX idx_links_click_count (click_count),
    INDEX idx_links_is_active (is_active),
    INDEX idx_links_tags USING GIN (tags),
    
    -- Constraints
    CONSTRAINT chk_short_code_length CHECK (LENGTH(short_code) >= 3),
    CONSTRAINT chk_short_code_alphanumeric CHECK (short_code ~ '^[a-zA-Z0-9_-]+$')
);

CREATE TRIGGER update_links_updated_at 
    BEFORE UPDATE ON links 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
```

### domains Table
```sql
CREATE TABLE domains (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    domain_name VARCHAR(255) UNIQUE NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    is_verified BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    verification_token VARCHAR(64),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    verified_at TIMESTAMP WITH TIME ZONE,
    
    -- Indexes
    INDEX idx_domains_domain_name (domain_name),
    INDEX idx_domains_user_id (user_id),
    INDEX idx_domains_is_verified (is_verified),
    
    -- Constraints
    CONSTRAINT chk_domain_name_format CHECK (domain_name ~ '^[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$')
);
```

### api_keys Table
```sql
CREATE TABLE api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    key_hash VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    last_used_at TIMESTAMP WITH TIME ZONE,
    expires_at TIMESTAMP WITH TIME ZONE,
    rate_limit INTEGER DEFAULT 1000,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Indexes
    INDEX idx_api_keys_user_id (user_id),
    INDEX idx_api_keys_key_hash (key_hash),
    INDEX idx_api_keys_expires_at (expires_at)
);
```

### click_events Table
```sql
CREATE TABLE click_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    link_id UUID NOT NULL REFERENCES links(id) ON DELETE CASCADE,
    ip_address INET,
    user_agent TEXT,
    referrer TEXT,
    country_code CHAR(2),
    city VARCHAR(100),
    device_type VARCHAR(20),  -- mobile, desktop, tablet
    browser VARCHAR(50),
    os VARCHAR(50),
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Indexes for analytics queries
    INDEX idx_click_events_link_id (link_id),
    INDEX idx_click_events_timestamp (timestamp),
    INDEX idx_click_events_country_code (country_code),
    INDEX idx_click_events_device_type (device_type),
    INDEX idx_click_events_link_id_timestamp (link_id, timestamp DESC)
    
    -- Partitioning consideration for large datasets:
    -- Could be partitioned by date (monthly) if expecting billions of rows
);

-- For anonymized analytics, we might want to hash IP addresses
-- This could be done in application layer or via trigger
```

### sessions Table (Optional)
```sql
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) UNIQUE NOT NULL,
    device_info TEXT,
    ip_address INET,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_used_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Indexes
    INDEX idx_sessions_user_id (user_id),
    INDEX idx_sessions_token_hash (token_hash),
    INDEX idx_sessions_expires_at (expires_at)
);
```

## Migration Files

### Initial Migration (`001_init_schema.up.sql`)
```sql
-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create update_updated_at_column function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    is_active BOOLEAN DEFAULT TRUE,
    is_admin BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_login_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_users_is_active ON users(is_active);

CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Create domains table
CREATE TABLE domains (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    domain_name VARCHAR(255) UNIQUE NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    is_verified BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    verification_token VARCHAR(64),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    verified_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_domains_domain_name ON domains(domain_name);
CREATE INDEX idx_domains_user_id ON domains(user_id);
CREATE INDEX idx_domains_is_verified ON domains(is_verified);

-- Create links table
CREATE TABLE links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    short_code VARCHAR(32) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    domain_id UUID REFERENCES domains(id) ON DELETE SET NULL,
    title VARCHAR(255),
    description TEXT,
    tags TEXT[] DEFAULT '{}',
    is_active BOOLEAN DEFAULT TRUE,
    expires_at TIMESTAMP WITH TIME ZONE,
    click_count BIGINT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_clicked_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_links_short_code ON links(short_code);
CREATE INDEX idx_links_user_id ON links(user_id);
CREATE INDEX idx_links_domain_id ON links(domain_id);
CREATE INDEX idx_links_created_at ON links(created_at);
CREATE INDEX idx_links_click_count ON links(click_count);
CREATE INDEX idx_links_is_active ON links(is_active);
CREATE INDEX idx_links_tags ON links USING GIN(tags);

CREATE TRIGGER update_links_updated_at 
    BEFORE UPDATE ON links 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Create api_keys table
CREATE TABLE api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    key_hash VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    last_used_at TIMESTAMP WITH TIME ZONE,
    expires_at TIMESTAMP WITH TIME ZONE,
    rate_limit INTEGER DEFAULT 1000,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_api_keys_user_id ON api_keys(user_id);
CREATE INDEX idx_api_keys_key_hash ON api_keys(key_hash);
CREATE INDEX idx_api_keys_expires_at ON api_keys(expires_at);

-- Create click_events table
CREATE TABLE click_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    link_id UUID NOT NULL REFERENCES links(id) ON DELETE CASCADE,
    ip_address INET,
    user_agent TEXT,
    referrer TEXT,
    country_code CHAR(2),
    city VARCHAR(100),
    device_type VARCHAR(20),
    browser VARCHAR(50),
    os VARCHAR(50),
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_click_events_link_id ON click_events(link_id);
CREATE INDEX idx_click_events_timestamp ON click_events(timestamp);
CREATE INDEX idx_click_events_country_code ON click_events(country_code);
CREATE INDEX idx_click_events_device_type ON click_events(device_type);
CREATE INDEX idx_click_events_link_id_timestamp ON click_events(link_id, timestamp DESC);
```

### Add Analytics Views (`002_add_analytics_views.up.sql`)
```sql
-- Create materialized view for daily link statistics
CREATE MATERIALIZED VIEW link_daily_stats AS
SELECT 
    link_id,
    DATE(timestamp) AS date,
    COUNT(*) AS clicks,
    COUNT(DISTINCT ip_address) AS unique_visitors,
    COUNT(DISTINCT country_code) AS countries_count,
    MODE() WITHIN GROUP (ORDER BY device_type) AS top_device,
    MODE() WITHIN GROUP (ORDER BY country_code) AS top_country
FROM click_events
WHERE timestamp >= NOW() - INTERVAL '90 days'
GROUP BY link_id, DATE(timestamp);

CREATE UNIQUE INDEX idx_link_daily_stats_link_id_date 
    ON link_daily_stats(link_id, date);

-- Create view for link summary statistics
CREATE VIEW link_summary_stats AS
SELECT 
    l.id AS link_id,
    l.short_code,
    l.click_count AS total_clicks,
    COUNT(DISTINCT ce.country_code) AS countries_reached,
    COUNT(DISTINCT DATE(ce.timestamp)) AS active_days,
    MIN(ce.timestamp) AS first_click_at,
    MAX(ce.timestamp) AS last_click_at,
    AVG(CASE WHEN ce.device_type = 'mobile' THEN 1 ELSE 0 END) * 100 AS mobile_percentage
FROM links l
LEFT JOIN click_events ce ON l.id = ce.link_id
GROUP BY l.id, l.short_code, l.click_count;

-- Create function to refresh materialized views
CREATE OR REPLACE FUNCTION refresh_analytics_views()
RETURNS VOID AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY link_daily_stats;
END;
$$ LANGUAGE plpgsql;
```

### Add Performance Indexes (`003_add_performance_indexes.up.sql`)
```sql
-- Add composite indexes for common query patterns
CREATE INDEX idx_links_user_active_created 
    ON links(user_id, is_active, created_at DESC);

CREATE INDEX idx_click_events_link_date_device 
    ON click_events(link_id, DATE(timestamp), device_type);

-- Add partial indexes for better performance
CREATE INDEX idx_active_links 
    ON links(short_code) 
    WHERE is_active = TRUE AND (expires_at IS NULL OR expires_at > NOW());

CREATE INDEX idx_unexpired_api_keys 
    ON api_keys(user_id, expires_at) 
    WHERE expires_at IS NULL OR expires_at > NOW();

-- Add index for tag searching
CREATE INDEX idx_links_tags_gin 
    ON links USING GIN(tags gin_trgm_ops);

-- Add index for URL pattern matching (if needed for admin searches)
CREATE INDEX idx_links_original_url_trgm 
    ON links USING GIN(original_url gin_trgm_ops);
```

### Add Admin Features (`004_add_admin_features.up.sql`)
```sql
-- Add audit logging table
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(50) NOT NULL,
    resource_type VARCHAR(50) NOT NULL,
    resource_id UUID,
    details JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX idx_audit_logs_resource ON audit_logs(resource_type, resource_id);

-- Add system settings table
CREATE TABLE system_settings (
    key VARCHAR(100) PRIMARY KEY,
    value TEXT,
    description TEXT,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL
);

INSERT INTO system_settings (key, value, description) VALUES
    ('max_links_per_user', '1000', 'Maximum number of links a user can create'),
    ('max_domains_per_user', '10', 'Maximum number of domains a user can add'),
    ('default_rate_limit', '1000', 'Default API rate limit per hour'),
    ('allow_custom_codes', 'true', 'Whether users can specify custom short codes'),
    ('require_email_verification', 'false', 'Whether email verification is required');

-- Add user quotas table
CREATE TABLE user_quotas (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    max_links INTEGER DEFAULT 1000,
    max_domains INTEGER DEFAULT 10,
    max_api_keys INTEGER DEFAULT 10,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

## Migration Rollback Files

### Initial Migration Rollback (`001_init_schema.down.sql`)
```sql
DROP TABLE IF EXISTS click_events CASCADE;
DROP TABLE IF EXISTS api_keys CASCADE;
DROP TABLE IF EXISTS links CASCADE;
DROP TABLE IF EXISTS domains CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP FUNCTION IF EXISTS update_updated_at_column() CASCADE;
```

## Database Configuration

### Connection Pool Settings
```yaml
# Recommended PostgreSQL configuration (postgresql.conf)
max_connections = 200
shared_buffers = 1GB
effective_cache_size = 3GB
maintenance_work_mem = 256MB
checkpoint_completion_target = 0.9
wal_buffers = 16MB
default_statistics_target = 100
random_page_cost = 1.1
effective_io_concurrency = 200
work_mem = 4MB
min_wal_size = 1GB
max_wal_size = 4GB
```

### Application Connection Settings
```go
// Recommended GORM configuration
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    PrepareStmt:            true,  // Use prepared statements
    SkipDefaultTransaction: true,  // Better performance
    Logger:                 logger.Default.LogMode(logger.Silent),
})

// Connection pool settings
sqlDB, err := db.DB()
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

## Data Retention Policies

### Click Events Retention
```sql
-- Create a retention policy for click events
-- Keep 90 days of detailed data, then aggregate and delete
CREATE OR REPLACE PROCEDURE cleanup_old_click_events()
LANGUAGE plpgsql
AS $$
BEGIN
    -- Archive old data first (optional)
    INSERT INTO click_events_archive
    SELECT * FROM click_events 
    WHERE timestamp < NOW() - INTERVAL '90 days';
    
    -- Delete old data
    DELETE FROM click_events 
    WHERE timestamp < NOW() - INTERVAL '90 days';
    
    -- Update link aggregates
    UPDATE links l
    SET click_count = (
        SELECT COUNT(*) 
        FROM click_events ce 
        WHERE ce.link_id = l.id
    )
    WHERE l.updated_at < NOW() - INTERVAL '1 day';
END;
$$;

-- Schedule with pg_cron if available
-- SELECT cron.schedule('cleanup-click-events', '0 2 * * *', 'CALL cleanup_old_click_events()');
```

## Backup Strategy

### Daily Backups
```bash
# Example backup script
pg_dump -U urlshortener -h localhost -d url_shortener \
  --format=custom --blobs --verbose --file=/backups/url_shortener_$(date +%Y%m%d).dump

# Backup only schema
pg_dump -U urlshortener -h localhost -d url_shortener \
  --schema-only --file=/backups/schema_$(date +%Y%m%d).sql
```

### Point-in-Time Recovery
Enable WAL archiving in `postgresql.conf`:
```ini
wal_level = replica
archive_mode = on
archive_command = 'cp %p /var/lib/postgresql/wal_archive/%f'
```

## Performance Optimization

### Query Optimization Examples

1. **Redirect Query** (most frequent):
```sql
-- Optimized for redirects
SELECT original_url, is_active, expires_at 
FROM links 
WHERE short_code = $1 
AND is_active = TRUE 
AND (expires_at IS NULL OR expires_at > NOW())
LIMIT 1;
```

2. **User Links with Pagination**:
```sql
SELECT * FROM links 
WHERE user_id = $1 
ORDER BY created_at DESC 
LIMIT $2 OFFSET $3;
```

3. **Analytics Aggregation**:
```sql
-- Use materialized view for daily stats
SELECT * FROM link_daily_stats 
WHERE link_id = $1 
AND date >= $2 
ORDER BY date DESC;
```

### Indexing Strategy

| Table | Index | Purpose | Size Estimate |
|-------|-------|---------|---------------|
| links | (short_code) | Redirect lookups | Small |
| links | (user_id, created_at) | User dashboard | Medium |
| links | (is_active, expires_at) | Cleanup jobs | Small |
| click_events | (link_id, timestamp) | Time-series analytics | Large |
| click_events | (timestamp) | Retention policies | Large |
| users | (email) | Authentication | Small |

## Security Considerations

### Row-Level Security (Optional)
```sql
-- Enable RLS
ALTER TABLE links ENABLE ROW LEVEL SECURITY;

-- Policy for users to see only their own links
CREATE POLICY user_links_policy ON links
    FOR ALL USING (user_id = current_user_id());

-- Function to get current user ID from JWT
CREATE OR REPLACE FUNCTION current_user_id()
RETURNS UUID AS $$
BEGIN
    RETURN NULLIF(current_setting('app.current_user_id', TRUE), '')::UUID;
END;
$$ LANGUAGE plpgsql;
```

### Data Encryption
- Passwords: bcrypt hashing in application layer
- API keys: HMAC-SHA256 hashing before storage
- Sensitive data: Consider PostgreSQL encryption or application-layer encryption

## Monitoring and Maintenance

### Key Metrics to Monitor
```sql
-- Database size
SELECT pg_size_pretty(pg_database_size('url_shortener'));

-- Table sizes
SELECT 
    table_name,
    pg_size_pretty(pg_total_relation_size(quote_ident(table_name))) as size
FROM information_schema.tables
WHERE table_schema = 'public'
ORDER BY pg_total_relation_size(quote_ident(table_name)) DESC;

-- Index usage statistics
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan as index_scans
FROM pg_stat_user_indexes
ORDER BY idx_scan DESC;
```

### Vacuum and Analyze
```sql
-- Manual vacuum for large tables
VACUUM ANALYZE click_events;

-- Set autovacuum parameters for large tables
ALTER TABLE click_events SET (
    autovacuum_vacuum_scale_factor = 0.05,
    autovacuum_analyze_scale_factor = 0.02
);
```

## Migration Management

### Using golang-migrate
```bash
# Create new migration
migrate create -ext sql -dir migrations -seq add_new_feature

# Apply migrations
migrate -path migrations -database "postgres://user:pass@localhost:5432/url_shortener?sslmode=disable" up

# Rollback last migration
migrate -path migrations -database "postgres://user:pass@localhost:5432/url_shortener?sslmode=disable" down 1

# Check migration status
migrate -path migrations -database "postgres://user:pass@localhost:5432/url_shortener?sslmode=disable" version
```

### Migration Best Practices
1. **Idempotent migrations**: Use `CREATE TABLE IF NOT EXISTS`
2. **Backward compatibility**: Maintain old columns during transition
3. **Data migrations**: Separate from schema migrations
4. **Testing**: Test migrations in staging before production
5. **Rollback plans**: Always have working down migrations

This schema provides a solid foundation for a scalable URL shortener service with support for analytics, user management, and custom domains.