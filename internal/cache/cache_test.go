package cache

import (
	"context"
	"testing"

	"github.com/MyNameIsWhaaat/shortener/internal/domain"
)

func TestNoOpCache(t *testing.T) {
	cache := &NoOpCache{}
	ctx := context.Background()

	// Test Get
	result, err := cache.Get(ctx, "test")
	if err != nil || result != nil {
		t.Errorf("Expected nil result and nil error, got %v, %v", result, err)
	}

	// Test Set
	url := &domain.URL{ShortCode: "test", OriginalURL: "http://example.com"}
	if err := cache.Set(ctx, "test", url); err != nil {
		t.Errorf("Set should not fail, got %v", err)
	}

	// Test Invalidate
	if err := cache.Invalidate(ctx, "test"); err != nil {
		t.Errorf("Invalidate should not fail, got %v", err)
	}

	// Test GetPopular
	results, err := cache.GetPopular(ctx, 10)
	if err != nil || results != nil {
		t.Errorf("GetPopular should return nil, got %v, %v", results, err)
	}

	// Test IncrementPopularity
	if err := cache.IncrementPopularity(ctx, "test"); err != nil {
		t.Errorf("IncrementPopularity should not fail, got %v", err)
	}

	// Test Close
	if err := cache.Close(); err != nil {
		t.Errorf("Close should not fail, got %v", err)
	}
}
