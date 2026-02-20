package cache

import (
	"context"

	"github.com/MyNameIsWhaaat/shortener/internal/domain"
)

type Cache interface {
	// Get retrieves a URL from cache
	Get(ctx context.Context, shortCode string) (*domain.URL, error)

	// Set stores a URL in cache with TTL
	Set(ctx context.Context, shortCode string, url *domain.URL) error

	// Invalidate removes a URL from cache
	Invalidate(ctx context.Context, shortCode string) error

	// GetPopular returns the most popular/recently accessed URLs
	GetPopular(ctx context.Context, limit int) ([]*domain.URL, error)

	// IncrementPopularity increments the popularity score for a URL
	IncrementPopularity(ctx context.Context, shortCode string) error

	// Close closes the cache connection
	Close() error
}

// NoOpCache is a cache implementation that does nothing (useful for testing)
type NoOpCache struct{}

func (n *NoOpCache) Get(ctx context.Context, shortCode string) (*domain.URL, error) {
	return nil, nil
}

func (n *NoOpCache) Set(ctx context.Context, shortCode string, url *domain.URL) error {
	return nil
}

func (n *NoOpCache) Invalidate(ctx context.Context, shortCode string) error {
	return nil
}

func (n *NoOpCache) GetPopular(ctx context.Context, limit int) ([]*domain.URL, error) {
	return nil, nil
}

func (n *NoOpCache) IncrementPopularity(ctx context.Context, shortCode string) error {
	return nil
}

func (n *NoOpCache) Close() error {
	return nil
}
