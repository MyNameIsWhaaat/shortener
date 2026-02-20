package service

import (
	"context"

	"github.com/MyNameIsWhaaat/shortener/internal/domain"
)

type ShortenerService interface {
    CreateShortURL(ctx context.Context, req *domain.CreateURLRequest) (*domain.CreateURLResponse, error)
    GetOriginalURL(ctx context.Context, shortCode string) (*domain.URL, error)
    TrackClick(ctx context.Context, shortCode, userAgent, ip, referer string) error
    GetAllURLs(ctx context.Context, limit int) ([]*domain.URL, error)
}

type AnalyticsService interface {
    GetAnalytics(ctx context.Context, shortCode string) (*domain.AnalyticsResponse, error)
    GetDailyStats(ctx context.Context, shortCode string, days int) (map[string]int64, error)
    GetMonthlyStats(ctx context.Context, shortCode string, months int) (map[string]int64, error)
    GetDeviceStats(ctx context.Context, shortCode string) (map[string]int64, error)
    GetRecentClicks(ctx context.Context, shortCode string, limit int) ([]domain.ClickEvent, error)
}