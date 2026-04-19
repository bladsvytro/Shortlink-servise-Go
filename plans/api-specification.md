# URL Shortener API Specification

## Base Information
- **Base URL**: `https://api.example.com/api/v1`
- **Content-Type**: `application/json`
- **Authentication**: Bearer JWT or API Key

## Authentication

### Register User
**POST** `/auth/register`

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "name": "John Doe"
}
```

**Response (201 Created):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "name": "John Doe",
  "created_at": "2024-01-15T10:30:00Z"
}
```

### Login
**POST** `/auth/login`

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

**Response (200 OK):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "bearer",
  "expires_in": 900,
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "name": "John Doe"
  }
}
```

### Refresh Token
**POST** `/auth/refresh`

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response:** Same as login response with new tokens.

### Get Profile
**GET** `/auth/profile`

**Headers:**
```
Authorization: Bearer {access_token}
```

**Response (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "name": "John Doe",
  "is_active": true,
  "is_admin": false,
  "created_at": "2024-01-15T10:30:00Z",
  "last_login_at": "2024-01-16T14:20:00Z",
  "stats": {
    "total_links": 42,
    "total_clicks": 1250,
    "active_domains": 2
  }
}
```

## Link Management

### Create Link
**POST** `/links`

**Headers:**
```
Authorization: Bearer {access_token}
```

**Request Body:**
```json
{
  "original_url": "https://example.com/very/long/url/path",
  "custom_code": "mycode",  // optional
  "title": "Example Page",  // optional
  "description": "A description",  // optional
  "tags": ["example", "demo"],  // optional
  "domain_id": "d550e8400-e29b-41d4-a716-446655440001",  // optional
  "expires_at": "2024-12-31T23:59:59Z"  // optional
}
```

**Response (201 Created):**
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440000",
  "short_code": "mycode",
  "short_url": "https://short.example.com/mycode",
  "original_url": "https://example.com/very/long/url/path",
  "title": "Example Page",
  "description": "A description",
  "tags": ["example", "demo"],
  "click_count": 0,
  "is_active": true,
  "expires_at": "2024-12-31T23:59:59Z",
  "created_at": "2024-01-15T10:30:00Z",
  "qr_code_url": "https://api.example.com/qr/mycode"  // optional
}
```

### List Links
**GET** `/links`

**Query Parameters:**
- `page` (optional, default: 1)
- `limit` (optional, default: 20, max: 100)
- `search` (optional) - search in title, description, URL
- `tags` (optional) - filter by tags (comma-separated)
- `sort` (optional) - `created_at`, `-created_at`, `click_count`, `-click_count`
- `active_only` (optional, default: true) - show only active links

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": "660e8400-e29b-41d4-a716-446655440000",
      "short_code": "abc123",
      "short_url": "https://short.example.com/abc123",
      "original_url": "https://example.com/page",
      "title": "Example",
      "click_count": 42,
      "created_at": "2024-01-15T10:30:00Z",
      "last_clicked_at": "2024-01-16T14:20:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 150,
    "total_pages": 8
  }
}
```

### Get Link Details
**GET** `/links/{id}`

**Response (200 OK):**
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440000",
  "short_code": "abc123",
  "short_url": "https://short.example.com/abc123",
  "original_url": "https://example.com/very/long/url/path",
  "title": "Example Page",
  "description": "A description",
  "tags": ["example", "demo"],
  "click_count": 42,
  "is_active": true,
  "expires_at": "2024-12-31T23:59:59Z",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-16T14:20:00Z",
  "last_clicked_at": "2024-01-16T14:20:00Z",
  "domain": {
    "id": "d550e8400-e29b-41d4-a716-446655440001",
    "domain_name": "custom.example.com",
    "is_verified": true
  }
}
```

### Update Link
**PUT** `/links/{id}`

**Request Body:** (partial update allowed)
```json
{
  "title": "Updated Title",
  "description": "Updated description",
  "tags": ["updated", "tags"],
  "is_active": false
}
```

**Response:** Updated link object.

### Delete Link
**DELETE** `/links/{id}`

**Response (204 No Content)**

## Analytics

### Get Link Analytics
**GET** `/links/{id}/analytics`

**Query Parameters:**
- `period` (optional) - `day`, `week`, `month`, `year`, `all` (default: `month`)
- `from` (optional) - ISO timestamp
- `to` (optional) - ISO timestamp

**Response (200 OK):**
```json
{
  "link_id": "660e8400-e29b-41d4-a716-446655440000",
  "total_clicks": 1250,
  "unique_visitors": 850,
  "click_through_rate": 0.65,
  "time_series": [
    {
      "date": "2024-01-01",
      "clicks": 45,
      "unique_visitors": 32
    }
  ],
  "top_countries": [
    {
      "country": "US",
      "clicks": 450,
      "percentage": 36.0
    }
  ],
  "top_referrers": [
    {
      "referrer": "https://twitter.com",
      "clicks": 120,
      "percentage": 9.6
    }
  ],
  "device_distribution": {
    "mobile": 620,
    "desktop": 500,
    "tablet": 130
  },
  "browser_distribution": {
    "chrome": 800,
    "safari": 300,
    "firefox": 100,
    "other": 50
  }
}
```

### Get Click Events
**GET** `/links/{id}/clicks`

**Query Parameters:**
- `page` (optional, default: 1)
- `limit` (optional, default: 50, max: 200)
- `from` (optional) - ISO timestamp
- `to` (optional) - ISO timestamp
- `country` (optional) - filter by country code
- `device_type` (optional) - `mobile`, `desktop`, `tablet`

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": "770e8400-e29b-41d4-a716-446655440000",
      "timestamp": "2024-01-16T14:20:00Z",
      "ip_address": "192.168.1.1",
      "country_code": "US",
      "city": "New York",
      "device_type": "mobile",
      "browser": "Chrome",
      "os": "Android",
      "referrer": "https://twitter.com"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 50,
    "total": 1250
  }
}
```

## Domain Management

### Add Domain
**POST** `/domains`

**Request Body:**
```json
{
  "domain_name": "custom.example.com"
}
```

**Response (201 Created):**
```json
{
  "id": "d550e8400-e29b-41d4-a716-446655440001",
  "domain_name": "custom.example.com",
  "is_verified": false,
  "is_active": true,
  "created_at": "2024-01-15T10:30:00Z",
  "verification_token": "abc123def456",
  "verification_instructions": "Add TXT record: short-verify=abc123def456"
}
```

### List Domains
**GET** `/domains`

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": "d550e8400-e29b-41d4-a716-446655440001",
      "domain_name": "custom.example.com",
      "is_verified": true,
      "is_active": true,
      "created_at": "2024-01-15T10:30:00Z",
      "verified_at": "2024-01-16T14:20:00Z",
      "link_count": 15
    }
  ]
}
```

### Verify Domain
**POST** `/domains/{id}/verify`

**Response (200 OK):**
```json
{
  "id": "d550e8400-e29b-41d4-a716-446655440001",
  "domain_name": "custom.example.com",
  "is_verified": true,
  "verified_at": "2024-01-16T14:20:00Z",
  "message": "Domain verified successfully"
}
```

### Delete Domain
**DELETE** `/domains/{id}`

**Response (204 No Content)**

## API Key Management

### Create API Key
**POST** `/apikeys`

**Request Body:**
```json
{
  "name": "Production API Key",
  "rate_limit": 1000,
  "expires_at": "2024-12-31T23:59:59Z"  // optional
}
```

**Response (201 Created):**
```json
{
  "id": "a550e8400-e29b-41d4-a716-446655440001",
  "name": "Production API Key",
  "api_key": "sk_live_abc123def456ghi789",  // Only shown once!
  "key_prefix": "sk_live_abc",
  "rate_limit": 1000,
  "expires_at": "2024-12-31T23:59:59Z",
  "created_at": "2024-01-15T10:30:00Z",
  "warning": "Store this API key securely. It will not be shown again."
}
```

### List API Keys
**GET** `/apikeys`

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": "a550e8400-e29b-41d4-a716-446655440001",
      "name": "Production API Key",
      "key_prefix": "sk_live_abc",
      "rate_limit": 1000,
      "last_used_at": "2024-01-16T14:20:00Z",
      "expires_at": "2024-12-31T23:59:59Z",
      "created_at": "2024-01-15T10:30:00Z"
    }
  ]
}
```

### Revoke API Key
**DELETE** `/apikeys/{id}`

**Response (204 No Content)**

## Public Redirect Endpoints

### Redirect (Default Domain)
**GET** `/{shortCode}`

**Response:**
- **302 Found** with `Location` header to original URL
- **404 Not Found** if link doesn't exist or is inactive
- **410 Gone** if link has expired

**Headers in response:**
- `Cache-Control: public, max-age=3600` (for active links)
- `X-Redirect-By: URL Shortener`

### Redirect (Custom Domain)
**GET** `/{shortCode}` on custom domain

Same behavior as above.

### Get Link Info (Public)
**GET** `/i/{shortCode}`

**Response (200 OK):**
```json
{
  "short_code": "abc123",
  "original_url": "https://example.com/very/long/url/path",
  "title": "Example Page",
  "description": "A description",
  "created_at": "2024-01-15T10:30:00Z",
  "click_count": 42
}
```

## Admin Endpoints

### List All Users
**GET** `/admin/users`

**Query Parameters:**
- `page`, `limit`, `search`, `active_only`

**Headers:**
```
Authorization: Bearer {admin_token}
```

**Response:** Paginated list of users with extended info.

### Toggle User Status
**PUT** `/admin/users/{id}/status`

**Request Body:**
```json
{
  "is_active": false,
  "reason": "Violation of terms"  // optional
}
```

### List All Links
**GET** `/admin/links`

**Query Parameters:** Extended filters available

### System Statistics
**GET** `/admin/stats`

**Response (200 OK):**
```json
{
  "total_users": 1250,
  "active_users": 980,
  "total_links": 42500,
  "active_links": 38900,
  "total_clicks": 1250000,
  "clicks_today": 1250,
  "clicks_this_week": 8750,
  "top_users": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "user@example.com",
      "link_count": 150,
      "click_count": 12500
    }
  ],
  "system_health": {
    "database": "healthy",
    "cache": "healthy",
    "uptime": "99.95%"
  }
}
```

## Error Responses

### Standard Error Format
```json
{
  "error": {
    "code": "validation_error",
    "message": "Invalid input parameters",
    "details": {
      "email": ["must be a valid email address"],
      "password": ["must be at least 8 characters"]
    },
    "request_id": "req_abc123def456"
  }
}
```

### Common Error Codes
- `authentication_required` - 401
- `invalid_credentials` - 401
- `insufficient_permissions` - 403
- `resource_not_found` - 404
- `validation_error` - 422
- `rate_limit_exceeded` - 429
- `internal_server_error` - 500

## Rate Limiting

### Limits
- **Authenticated users**: 1000 requests/hour
- **API keys**: Configurable per key (default: 1000/hour)
- **Public endpoints**: 100 requests/hour per IP
- **Admin endpoints**: 500 requests/hour

### Headers
```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 950
X-RateLimit-Reset: 1612345678
```

## Pagination

All list endpoints support pagination with consistent format:
- `page`: Current page (1-indexed)
- `limit`: Items per page (default 20, max 100)
- `total`: Total items
- `total_pages`: Total pages

## Sorting

Supported sort fields:
- `created_at` (default) or `-created_at`
- `updated_at` or `-updated_at`
- `click_count` or `-click_count`
- `last_clicked_at` or `-last_clicked_at`

## Filtering

Common filter parameters:
- `active_only` (boolean)
- `from` / `to` (timestamp)
- `search` (text search)
- `tags` (comma-separated)

## Webhooks (Optional)

### Events
- `link.created`
- `link.updated`
- `link.deleted`
- `click.registered`
- `domain.verified`

### Webhook Configuration
**POST** `/webhooks`
**DELETE** `/webhooks/{id}`

## Bulk Operations (Optional)

### Bulk Create Links
**POST** `/links/bulk`

**Request Body:**
```json
{
  "links": [
    {
      "original_url": "https://example.com/1",
      "custom_code": "link1"
    },
    {
      "original_url": "https://example.com/2"
    }
  ]
}
```

**Response:** Array of created links with status for each.

### Bulk Update Links
**PUT** `/links/bulk`

## API Versioning

- Current version: `v1`
- Version in URL path: `/api/v1/`
- Deprecation notice via `Deprecation` header
- Sunset policy: 6 months notice for breaking changes