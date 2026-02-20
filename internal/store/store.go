package store

import (
	"context"

	"github.com/MyNameIsWhaaat/shortener/internal/domain"
)

type URLStore interface {
	CreateURL(ctx context.Context, url *domain.URL) error
	GetURLByShortCode(ctx context.Context, shortCode string) (*domain.URL, error)
	IncrementClicks(ctx context.Context, shortCode string) error
	CheckShortCodeExists(ctx context.Context, shortCode string) (bool, error)
	GetAllURLs(ctx context.Context, limit int) ([]*domain.URL, error)
}

type AnalyticsStore interface {
	SaveClickEvent(ctx context.Context, event *domain.ClickEvent) error
	GetDailyStats(ctx context.Context, shortCode string, days int) (map[string]int64, error)
	GetMonthlyStats(ctx context.Context, shortCode string, months int) (map[string]int64, error)
	GetDeviceStats(ctx context.Context, shortCode string) (map[string]int64, error)
	GetRecentClicks(ctx context.Context, shortCode string, limit int) ([]domain.ClickEvent, error)
	GetAnalytics(ctx context.Context, shortCode string) (*domain.AnalyticsResponse, error)
}

type Store interface {
	URLStore
	AnalyticsStore
	Ping(ctx context.Context) error
	Close() error
}
