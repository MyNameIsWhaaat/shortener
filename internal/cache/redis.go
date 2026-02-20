package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/MyNameIsWhaaat/shortener/internal/domain"
	"github.com/redis/go-redis/v9"
)

const (
	urlKeyPrefix       = "url:"
	popularityKeyPrefix = "pop:"
	popularitySetKey   = "urls:popular"
	defaultTTL         = 24 * time.Hour
)

type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisCache(redisAddr, password string, db int, ttl time.Duration) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	if ttl == 0 {
		ttl = defaultTTL
	}

	log.Printf("Redis cache connected successfully (TTL: %v)", ttl)
	return &RedisCache{
		client: client,
		ttl:    ttl,
	}, nil
}

// Get retrieves a URL from cache
func (rc *RedisCache) Get(ctx context.Context, shortCode string) (*domain.URL, error) {
	key := urlKeyPrefix + shortCode
	val, err := rc.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	}
	if err != nil {
		log.Printf("Redis get error: %v", err)
		return nil, nil // Don't fail on cache error, allow fallback
	}

	var url domain.URL
	if err := json.Unmarshal([]byte(val), &url); err != nil {
		log.Printf("Failed to unmarshal cached URL: %v", err)
		return nil, nil
	}

	// Increment popularity when accessed
	go func() {
		if err := rc.IncrementPopularity(context.Background(), shortCode); err != nil {
			log.Printf("Failed to increment popularity: %v", err)
		}
	}()

	log.Printf("Cache hit for %s", shortCode)
	return &url, nil
}

// Set stores a URL in cache with TTL
func (rc *RedisCache) Set(ctx context.Context, shortCode string, url *domain.URL) error {
	key := urlKeyPrefix + shortCode
	data, err := json.Marshal(url)
	if err != nil {
		return fmt.Errorf("failed to marshal URL: %w", err)
	}

	if err := rc.client.Set(ctx, key, string(data), rc.ttl).Err(); err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	log.Printf("Cached URL %s (TTL: %v)", shortCode, rc.ttl)
	return nil
}

// Invalidate removes a URL from cache
func (rc *RedisCache) Invalidate(ctx context.Context, shortCode string) error {
	key := urlKeyPrefix + shortCode
	popKey := popularityKeyPrefix + shortCode

	if err := rc.client.Del(ctx, key, popKey).Err(); err != nil {
		return fmt.Errorf("failed to invalidate cache: %w", err)
	}

	// Remove from sorted set
	if err := rc.client.ZRem(ctx, popularitySetKey, shortCode).Err(); err != nil {
		log.Printf("Failed to remove from popularity set: %v", err)
	}

	log.Printf("Invalidated cache for %s", shortCode)
	return nil
}

// GetPopular returns the most popular/recently accessed URLs
func (rc *RedisCache) GetPopular(ctx context.Context, limit int) ([]*domain.URL, error) {
	if limit <= 0 {
		limit = 10
	}

	// Get top N short codes by popularity score
	codes, err := rc.client.ZRevRange(ctx, popularitySetKey, 0, int64(limit-1)).Result()
	if err != nil {
		log.Printf("Failed to get popular codes: %v", err)
		return nil, nil // Don't fail on cache error
	}

	if len(codes) == 0 {
		return nil, nil
	}

	urls := make([]*domain.URL, 0, len(codes))
	for _, code := range codes {
		url, err := rc.Get(ctx, code)
		if err != nil {
			log.Printf("Failed to get popular URL %s: %v", code, err)
			continue
		}
		if url != nil {
			urls = append(urls, url)
		}
	}

	return urls, nil
}

// IncrementPopularity increments the popularity score for a URL
func (rc *RedisCache) IncrementPopularity(ctx context.Context, shortCode string) error {
	// Increment score in sorted set
	if err := rc.client.ZIncrBy(ctx, popularitySetKey, 1, shortCode).Err(); err != nil {
		log.Printf("Failed to increment popularity: %v", err)
		return nil // Don't fail on cache error
	}

	return nil
}

// Close closes the Redis connection
func (rc *RedisCache) Close() error {
	return rc.client.Close()
}
