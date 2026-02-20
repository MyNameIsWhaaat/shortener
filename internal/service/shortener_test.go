package service

import (
	"context"
	"testing"

	"github.com/MyNameIsWhaaat/shortener/internal/cache"
	"github.com/MyNameIsWhaaat/shortener/internal/domain"
)

// MockURLStore for testing
type MockURLStore struct {
	urls map[string]*domain.URL
}

func NewMockURLStore() *MockURLStore {
	return &MockURLStore{
		urls: make(map[string]*domain.URL),
	}
}

func (m *MockURLStore) CreateURL(ctx context.Context, url *domain.URL) error {
	m.urls[url.ShortCode] = url
	return nil
}

func (m *MockURLStore) GetURLByShortCode(ctx context.Context, shortCode string) (*domain.URL, error) {
	if url, ok := m.urls[shortCode]; ok {
		return url, nil
	}
	return nil, ErrURLNotFound
}

func (m *MockURLStore) IncrementClicks(ctx context.Context, shortCode string) error {
	if url, ok := m.urls[shortCode]; ok {
		url.Clicks++
		return nil
	}
	return ErrURLNotFound
}

func (m *MockURLStore) CheckShortCodeExists(ctx context.Context, shortCode string) (bool, error) {
	_, exists := m.urls[shortCode]
	return exists, nil
}

func (m *MockURLStore) GetAllURLs(ctx context.Context, limit int) ([]*domain.URL, error) {
	urls := make([]*domain.URL, 0)
	for _, url := range m.urls {
		urls = append(urls, url)
	}
	return urls, nil
}

// MockAnalyticsStore for testing
type MockAnalyticsStore struct {
	events map[string][]domain.ClickEvent
}

func NewMockAnalyticsStore() *MockAnalyticsStore {
	return &MockAnalyticsStore{
		events: make(map[string][]domain.ClickEvent),
	}
}

func (m *MockAnalyticsStore) SaveClickEvent(ctx context.Context, event *domain.ClickEvent) error {
	m.events[event.ShortCode] = append(m.events[event.ShortCode], *event)
	return nil
}

func (m *MockAnalyticsStore) GetDailyStats(ctx context.Context, shortCode string, days int) (map[string]int64, error) {
	return map[string]int64{}, nil
}

func (m *MockAnalyticsStore) GetMonthlyStats(ctx context.Context, shortCode string, months int) (map[string]int64, error) {
	return map[string]int64{}, nil
}

func (m *MockAnalyticsStore) GetDeviceStats(ctx context.Context, shortCode string) (map[string]int64, error) {
	return map[string]int64{}, nil
}

func (m *MockAnalyticsStore) GetRecentClicks(ctx context.Context, shortCode string, limit int) ([]domain.ClickEvent, error) {
	if events, ok := m.events[shortCode]; ok {
		if len(events) > limit {
			return events[:limit], nil
		}
		return events, nil
	}
	return []domain.ClickEvent{}, nil
}

func (m *MockAnalyticsStore) GetAnalytics(ctx context.Context, shortCode string) (*domain.AnalyticsResponse, error) {
	return &domain.AnalyticsResponse{
		ShortCode:   shortCode,
		TotalClicks: 0,
	}, nil
}

// Tests
func TestCreateShortURL(t *testing.T) {
	urlStore := NewMockURLStore()
	analyticsStore := NewMockAnalyticsStore()
	cacheClient := &cache.NoOpCache{}

	service := NewShortenerService(urlStore, "http://localhost:8080", 6, analyticsStore, cacheClient)

	tests := []struct {
		name    string
		req     *domain.CreateURLRequest
		wantErr bool
	}{
		{
			name:    "valid URL",
			req:     &domain.CreateURLRequest{URL: "https://example.com"},
			wantErr: false,
		},
		{
			name:    "empty URL",
			req:     &domain.CreateURLRequest{URL: ""},
			wantErr: true,
		},
		{
			name:    "invalid URL (no scheme)",
			req:     &domain.CreateURLRequest{URL: "example.com"},
			wantErr: true,
		},
		{
			name:    "valid custom alias",
			req:     &domain.CreateURLRequest{URL: "https://example.com", CustomAlias: stringPtr("myalias")},
			wantErr: false,
		},
		{
			name:    "invalid custom alias (special chars)",
			req:     &domain.CreateURLRequest{URL: "https://example.com", CustomAlias: stringPtr("my@alias")},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.CreateShortURL(context.Background(), tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if !tt.wantErr && resp == nil {
				t.Error("expected response, got nil")
			}
		})
	}
}

func TestGetOriginalURL(t *testing.T) {
	urlStore := NewMockURLStore()
	analyticsStore := NewMockAnalyticsStore()
	cacheClient := &cache.NoOpCache{}

	service := NewShortenerService(urlStore, "http://localhost:8080", 6, analyticsStore, cacheClient)

	// Add a URL
	testURL := &domain.URL{
		ShortCode:   "abc123",
		OriginalURL: "https://example.com",
		Clicks:      0,
	}
	if err := urlStore.CreateURL(context.Background(), testURL); err != nil {
		t.Fatalf("Failed to create URL: %v", err)
	}

	tests := []struct {
		name      string
		shortCode string
		wantErr   bool
	}{
		{
			name:      "existing short code",
			shortCode: "abc123",
			wantErr:   false,
		},
		{
			name:      "non-existing short code",
			shortCode: "xyz789",
			wantErr:   true,
		},
		{
			name:      "empty short code",
			shortCode: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.GetOriginalURL(context.Background(), tt.shortCode)

			if (err != nil) != tt.wantErr {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if !tt.wantErr && result == nil {
				t.Error("expected result, got nil")
			}
		})
	}
}

func TestTrackClick(t *testing.T) {
	urlStore := NewMockURLStore()
	analyticsStore := NewMockAnalyticsStore()
	cacheClient := &cache.NoOpCache{}

	service := NewShortenerService(urlStore, "http://localhost:8080", 6, analyticsStore, cacheClient)

	// Add a URL
	testURL := &domain.URL{
		ShortCode:   "abc123",
		OriginalURL: "https://example.com",
		Clicks:      0,
	}
	if err := urlStore.CreateURL(context.Background(), testURL); err != nil {
		t.Fatalf("Failed to create URL: %v", err)
	}

	// Track a click
	err := service.TrackClick(context.Background(), "abc123", "Mozilla/5.0", "192.168.1.1", "https://google.com")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify clicks incremented
	url, _ := urlStore.GetURLByShortCode(context.Background(), "abc123")
	if url.Clicks != 1 {
		t.Errorf("expected 1 click, got %d", url.Clicks)
	}

	// Verify event recorded
	events, _ := analyticsStore.GetRecentClicks(context.Background(), "abc123", 10)
	if len(events) != 1 {
		t.Errorf("expected 1 event, got %d", len(events))
	}
}

func TestGetAllURLs(t *testing.T) {
	urlStore := NewMockURLStore()
	analyticsStore := NewMockAnalyticsStore()
	cacheClient := &cache.NoOpCache{}

	service := NewShortenerService(urlStore, "http://localhost:8080", 6, analyticsStore, cacheClient)

	// Add multiple URLs
	for i := 1; i <= 3; i++ {
		url := &domain.URL{
			ShortCode:   "code" + string(rune('0'+i)),
			OriginalURL: "https://example.com/" + string(rune('0'+i)),
			Clicks:      0,
		}
		if err := urlStore.CreateURL(context.Background(), url); err != nil {
			t.Fatalf("Failed to create URL: %v", err)
		}
	}

	urls, err := service.GetAllURLs(context.Background(), 10)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(urls) != 3 {
		t.Errorf("expected 3 URLs, got %d", len(urls))
	}
}

func TestValidateURL(t *testing.T) {
	urlStore := NewMockURLStore()
	analyticsStore := NewMockAnalyticsStore()
	cacheClient := &cache.NoOpCache{}

	service := NewShortenerService(urlStore, "http://localhost:8080", 6, analyticsStore, cacheClient).(*shortenerService)

	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"valid http URL", "http://example.com", false},
		{"valid https URL", "https://example.com/path", false},
		{"empty URL", "", true},
		{"no scheme", "example.com", true},
		{"invalid scheme", "ftp://example.com", true},
		{"no host", "https://", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateShortCode(t *testing.T) {
	urlStore := NewMockURLStore()
	analyticsStore := NewMockAnalyticsStore()
	cacheClient := &cache.NoOpCache{}

	service := NewShortenerService(urlStore, "http://localhost:8080", 6, analyticsStore, cacheClient).(*shortenerService)

	tests := []struct {
		name    string
		code    string
		wantErr bool
	}{
		{"valid code", "abc123", false},
		{"valid code with underscore", "abc_123", false},
		{"valid code with dash", "abc-123", false},
		{"empty code", "", true},
		{"too long code", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz123", true},
		{"invalid special char", "abc@123", true},
		{"invalid space", "abc 123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateShortCode(tt.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// Helper function
func stringPtr(s string) *string {
	return &s
}
