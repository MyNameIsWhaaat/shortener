# Redis Caching Implementation Summary

## Overview

Implemented Redis-based caching system for popular links with automatic popularity tracking and TTL-based expiration.

## Components Implemented

### 1. Cache Interface (`internal/cache/cache.go`)

- **Cache Interface**: Defines operations for URL caching
  - `Get(ctx, shortCode)` - Retrieve cached URL
  - `Set(ctx, shortCode, url)` - Store URL with TTL
  - `Invalidate(ctx, shortCode)` - Remove from cache
  - `GetPopular(ctx, limit)` - Get trending URLs
  - `IncrementPopularity(ctx, shortCode)` - Track access frequency
  - `Close()` - Cleanup resources

- **NoOpCache**: Fallback implementation for when Redis is unavailable

### 2. Redis Implementation (`internal/cache/redis.go`)

- **RedisCache struct**: Manages Redis connection and operations
- **Key Prefixes**:
  - `url:` - Stores serialized URL objects (JSON)
  - `pop:` - Tracks popularity scores
  - `urls:popular` - Sorted set of popular short codes

- **Features**:
  - Automatic TTL management (default 24 hours, configurable)
  - Popularity tracking via Redis Sorted Set
  - Automatic background popularity increment on cache hits
  - Graceful degradation if Redis unavailable

### 3. Service Integration (`internal/service/shortener.go`)

- **Cache Field**: Added to `shortenerService`
- **Updated Methods**:
  - `CreateShortURL()` - Caches newly created URLs
  - `GetOriginalURL()` - Cache-first strategy with DB fallback
  - _NEW_ `GetPopularURLs()` - Returns trending URLs from cache

- **Caching Strategy**:
  1. Check Redis cache first on redirect
  2. If miss, query PostgreSQL
  3. Auto-cache retrieved URLs
  4. Population score increases on each access

### 4. API Handler (`internal/httpapi/handler.go`)

- _NEW_ `GetPopularURLs()` endpoint handler
- Query parameter: `limit` (default 10, max 100)
- Returns sorted list of most-accessed URLs

### 5. Router Configuration (`internal/api/router.go`)

- _NEW_ Route: `GET /api/urls/popular?limit=10`

### 6. Main Application (`cmd/shortener/main.go`)

- Initializes Redis cache at startup
- Graceful fallback to NoOpCache if Redis unavailable
- Proper connection cleanup with defer

### 7. Dependencies (`go.mod`)

- Added: `github.com/redis/go-redis/v9 v9.5.1`

## Configuration

Redis settings already available in `internal/config/config.go`:

```
REDIS_HOST=redis        (default)
REDIS_PORT=6379         (default)
REDIS_PASSWORD=         (default empty)
REDIS_DB=0              (default)
CACHE_TTL=24h           (default)
```

## Docker Integration

- Redis service in docker-compose.yaml (already configured)
- Image: redis:7-alpine
- Port: 6379
- Health check enabled

## Usage Examples

### 1. Shorten a URL

```bash
curl -X POST http://localhost:8080/api/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com/very/long/url"}'
```

→ URL cached automatically

### 2. Access shortened URL (Popular)

```bash
GET http://localhost:8080/s/abc123
```

→ Loaded from cache (if available), popularity incremented

### 3. Get Popular URLs

```bash
GET http://localhost:8080/api/urls/popular?limit=10
```

→ Returns top 10 most-accessed URLs (from Redis Sorted Set)

## Performance Benefits

- **Reduced DB Load**: Frequently accessed URLs served from cache
- **Lower Latency**: Cache hits are ~100x faster than DB queries
- **Automatic Tracking**: Popularity data collected without explicit analytics calls
- **Configurable TTL**: Adjust cache lifetime per deployment needs

## Resilience

- **Graceful Degradation**: Application works without Redis
- **Error Logging**: Redis failures logged but don't crash service
- **Async Operations**: Popularity tracking runs async, non-blocking

## Key Files Modified

1. `cmd/shortener/main.go` - Initialize cache
2. `internal/service/shortener.go` - Implement cache strategy
3. `internal/service/service.go` - Add interface method
4. `internal/httpapi/handler.go` - Add handler
5. `internal/api/router.go` - Add route
6. `go.mod` - Add redis dependency
7. `internal/httpapi/redirect.go` - Add timeout context

## Files Created

1. `internal/cache/cache.go` - Cache interface
2. `internal/cache/redis.go` - Redis implementation
