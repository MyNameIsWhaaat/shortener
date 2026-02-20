package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MyNameIsWhaaat/shortener/internal/cache"
	"github.com/MyNameIsWhaaat/shortener/internal/domain"
	"github.com/MyNameIsWhaaat/shortener/internal/service"
)

// Setup for handler tests
func setupTestHandler(t *testing.T) *Handler {
	urlStore := &testURLStore{
		urls: make(map[string]*domain.URL),
	}
	analyticsStore := &testAnalyticsStore{
		events: make(map[string][]domain.ClickEvent),
	}
	cacheClient := &cache.NoOpCache{}

	shortenerService := service.NewShortenerService(
		urlStore,
		"http://localhost:8080",
		6,
		analyticsStore,
		cacheClient,
	)
	analyticsService := service.NewAnalyticsService(analyticsStore)

	return NewHandler(shortenerService, analyticsService)
}

type testURLStore struct {
	urls map[string]*domain.URL
}

func (t *testURLStore) CreateURL(ctx context.Context, url *domain.URL) error {
	t.urls[url.ShortCode] = url
	return nil
}

func (t *testURLStore) GetURLByShortCode(ctx context.Context, shortCode string) (*domain.URL, error) {
	if url, ok := t.urls[shortCode]; ok {
		return url, nil
	}
	return nil, service.ErrURLNotFound
}

func (t *testURLStore) IncrementClicks(ctx context.Context, shortCode string) error {
	if url, ok := t.urls[shortCode]; ok {
		url.Clicks++
		return nil
	}
	return service.ErrURLNotFound
}

func (t *testURLStore) CheckShortCodeExists(ctx context.Context, shortCode string) (bool, error) {
	_, exists := t.urls[shortCode]
	return exists, nil
}

func (t *testURLStore) GetAllURLs(ctx context.Context, limit int) ([]*domain.URL, error) {
	urls := make([]*domain.URL, 0)
	for _, url := range t.urls {
		urls = append(urls, url)
		if limit > 0 && len(urls) >= limit {
			break
		}
	}
	return urls, nil
}

type testAnalyticsStore struct {
	events map[string][]domain.ClickEvent
}

func (t *testAnalyticsStore) SaveClickEvent(ctx context.Context, event *domain.ClickEvent) error {
	t.events[event.ShortCode] = append(t.events[event.ShortCode], *event)
	return nil
}

func (t *testAnalyticsStore) GetDailyStats(ctx context.Context, shortCode string, days int) (map[string]int64, error) {
	return map[string]int64{}, nil
}

func (t *testAnalyticsStore) GetMonthlyStats(ctx context.Context, shortCode string, months int) (map[string]int64, error) {
	return map[string]int64{}, nil
}

func (t *testAnalyticsStore) GetDeviceStats(ctx context.Context, shortCode string) (map[string]int64, error) {
	return map[string]int64{}, nil
}

func (t *testAnalyticsStore) GetRecentClicks(ctx context.Context, shortCode string, limit int) ([]domain.ClickEvent, error) {
	if events, ok := t.events[shortCode]; ok {
		if len(events) > limit {
			return events[:limit], nil
		}
		return events, nil
	}
	return []domain.ClickEvent{}, nil
}

func (t *testAnalyticsStore) GetAnalytics(ctx context.Context, shortCode string) (*domain.AnalyticsResponse, error) {
	return &domain.AnalyticsResponse{
		ShortCode:   shortCode,
		TotalClicks: 0,
	}, nil
}

// Tests
func TestHealthHandler(t *testing.T) {
	handler := setupTestHandler(t)
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler.Health(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response map[string]string
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("expected status 'ok', got %s", response["status"])
	}
}

func TestShortenHandler(t *testing.T) {
	handler := setupTestHandler(t)

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		expectResponse bool
	}{
		{
			name: "valid request",
			requestBody: map[string]interface{}{
				"url": "https://example.com",
			},
			expectedStatus: http.StatusCreated,
			expectResponse: true,
		},
		{
			name: "invalid request (bad format)",
			requestBody: map[string]interface{}{
				"url": "not a url",
			},
			expectedStatus: http.StatusInternalServerError,
			expectResponse: false,
		},
		{
			name:           "invalid JSON",
			requestBody:    "invalid",
			expectedStatus: http.StatusBadRequest,
			expectResponse: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/shorten", bytes.NewReader(body))
			w := httptest.NewRecorder()

			handler.Shorten(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectResponse {
				var response map[string]string
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Errorf("failed to decode response: %v", err)
				}

				if response["short_code"] == "" {
					t.Error("expected short_code in response")
				}
			}
		})
	}
}

func TestGetAllURLsHandler(t *testing.T) {
	handler := setupTestHandler(t)

	req := httptest.NewRequest("GET", "/api/urls", nil)
	w := httptest.NewRecorder()

	// Create a URL first
	createReq := httptest.NewRequest("POST", "/api/shorten",
		bytes.NewReader([]byte(`{"url": "https://example.com"}`)),
	)
	createW := httptest.NewRecorder()
	handler.Shorten(createW, createReq)

	// Now get all URLs
	handler.GetAllURLs(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response []*domain.URL
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if len(response) < 1 {
		t.Error("expected at least 1 URL in response")
	}
}

func TestRespondError(t *testing.T) {
	handler := setupTestHandler(t)
	w := httptest.NewRecorder()

	handler.respondError(w, "test error", http.StatusBadRequest)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var response map[string]string
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["error"] != "test error" {
		t.Errorf("expected error message 'test error', got %s", response["error"])
	}
}
