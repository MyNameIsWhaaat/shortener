package store

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MyNameIsWhaaat/shortener/internal/domain"
)

var errNotFound = errors.New("not found")

// Mock implementations for testing interfaces
type mockURLStore struct {
	urls map[string]*domain.URL
}

func newMockURLStore() *mockURLStore {
	return &mockURLStore{
		urls: make(map[string]*domain.URL),
	}
}

func (m *mockURLStore) CreateURL(ctx context.Context, url *domain.URL) error {
	m.urls[url.ShortCode] = url
	return nil
}

func (m *mockURLStore) GetURLByShortCode(ctx context.Context, shortCode string) (*domain.URL, error) {
	if url, ok := m.urls[shortCode]; ok {
		return url, nil
	}
	return nil, errNotFound
}

func (m *mockURLStore) IncrementClicks(ctx context.Context, shortCode string) error {
	if url, ok := m.urls[shortCode]; ok {
		url.Clicks++
		return nil
	}
	return errNotFound
}

func (m *mockURLStore) CheckShortCodeExists(ctx context.Context, shortCode string) (bool, error) {
	_, exists := m.urls[shortCode]
	return exists, nil
}

func (m *mockURLStore) GetAllURLs(ctx context.Context, limit int) ([]*domain.URL, error) {
	urls := make([]*domain.URL, 0)
	count := 0
	for _, url := range m.urls {
		if limit > 0 && count >= limit {
			break
		}
		urls = append(urls, url)
		count++
	}
	return urls, nil
}

// Tests
func TestURLStoreInterface(t *testing.T) {
	store := newMockURLStore()
	ctx := context.Background()

	// Test CreateURL
	url := &domain.URL{
		ShortCode:   "test123",
		OriginalURL: "https://example.com",
		CreatedAt:   time.Now(),
		Clicks:      0,
	}

	err := store.CreateURL(ctx, url)
	if err != nil {
		t.Errorf("CreateURL failed: %v", err)
	}

	// Test GetURLByShortCode
	retrieved, err := store.GetURLByShortCode(ctx, "test123")
	if err != nil {
		t.Errorf("GetURLByShortCode failed: %v", err)
	}

	if retrieved.ShortCode != "test123" {
		t.Errorf("expected ShortCode test123, got %s", retrieved.ShortCode)
	}

	// Test CheckShortCodeExists
	exists, err := store.CheckShortCodeExists(ctx, "test123")
	if err != nil || !exists {
		t.Error("CheckShortCodeExists should return true for existing code")
	}

	notExists, err := store.CheckShortCodeExists(ctx, "notexist")
	if err != nil || notExists {
		t.Error("CheckShortCodeExists should return false for non-existing code")
	}

	// Test IncrementClicks
	err = store.IncrementClicks(ctx, "test123")
	if err != nil {
		t.Errorf("IncrementClicks failed: %v", err)
	}

	retrieved, _ = store.GetURLByShortCode(ctx, "test123")
	if retrieved.Clicks != 1 {
		t.Errorf("expected 1 click, got %d", retrieved.Clicks)
	}

	// Test GetAllURLs
	url2 := &domain.URL{
		ShortCode:   "test456",
		OriginalURL: "https://google.com",
		CreatedAt:   time.Now(),
		Clicks:      0,
	}
	if err := store.CreateURL(ctx, url2); err != nil {
		t.Fatalf("Failed to create URL: %v", err)
	}

	urls, err := store.GetAllURLs(ctx, 10)
	if err != nil {
		t.Errorf("GetAllURLs failed: %v", err)
	}

	if len(urls) != 2 {
		t.Errorf("expected 2 URLs, got %d", len(urls))
	}

	// Test GetAllURLs with limit
	urls, err = store.GetAllURLs(ctx, 1)
	if err != nil {
		t.Errorf("GetAllURLs with limit failed: %v", err)
	}
	if len(urls) != 1 {
		t.Errorf("expected 1 URL with limit=1, got %d", len(urls))
	}
}

func TestGetURLByShortCodeNotFound(t *testing.T) {
	store := newMockURLStore()
	ctx := context.Background()

	_, err := store.GetURLByShortCode(ctx, "nonexistent")
	if err == nil {
		t.Error("expected error for non-existent URL")
	}
}

func TestIncrementClicksNotFound(t *testing.T) {
	store := newMockURLStore()
	ctx := context.Background()

	err := store.IncrementClicks(ctx, "nonexistent")
	if err == nil {
		t.Error("expected error when incrementing non-existent URL")
	}
}
