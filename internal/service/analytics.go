package service

import (
	"context"
	"fmt"
	"log"

	"github.com/MyNameIsWhaaat/shortener/internal/domain"
	"github.com/MyNameIsWhaaat/shortener/internal/store"
)

type analyticsService struct {
    analyticsStore store.AnalyticsStore
}

func NewAnalyticsService(analyticsStore store.AnalyticsStore) AnalyticsService {
    return &analyticsService{
        analyticsStore: analyticsStore,
    }
}

func (s *analyticsService) GetAnalytics(ctx context.Context, shortCode string) (*domain.AnalyticsResponse, error) {
    if shortCode == "" {
        return nil, ErrInvalidShortCode
    }

    analytics, err := s.analyticsStore.GetAnalytics(ctx, shortCode)
    if err != nil {
        return nil, fmt.Errorf("failed to get analytics from store: %w", err)
    }

    return analytics, nil
}

func (s *analyticsService) GetDailyStats(ctx context.Context, shortCode string, days int) (map[string]int64, error) {
    log.Printf("Service GetDailyStats for %s, days=%d", shortCode, days)

    
    if shortCode == "" {
        return nil, ErrInvalidShortCode
    }

    if days <= 0 || days > 365 {
        days = 30
    }

    stats, err := s.analyticsStore.GetDailyStats(ctx, shortCode, days)
    if err != nil {
        return nil, fmt.Errorf("failed to get daily stats: %w", err)
    }

    return stats, nil
}

func (s *analyticsService) GetMonthlyStats(ctx context.Context, shortCode string, months int) (map[string]int64, error) {
    if shortCode == "" {
        return nil, ErrInvalidShortCode
    }

    if months <= 0 || months > 24 {
        months = 12
    }

    stats, err := s.analyticsStore.GetMonthlyStats(ctx, shortCode, months)
    if err != nil {
        return nil, fmt.Errorf("failed to get monthly stats: %w", err)
    }

    return stats, nil
}

func (s *analyticsService) GetDeviceStats(ctx context.Context, shortCode string) (map[string]int64, error) {
    if shortCode == "" {
        return nil, ErrInvalidShortCode
    }

    log.Printf("Service GetDeviceStats for %s", shortCode)

    stats, err := s.analyticsStore.GetDeviceStats(ctx, shortCode)
    if err != nil {
        return nil, fmt.Errorf("failed to get device stats: %w", err)
    }

    return stats, nil
}

func (s *analyticsService) GetRecentClicks(ctx context.Context, shortCode string, limit int) ([]domain.ClickEvent, error) {
    if shortCode == "" {
        return nil, ErrInvalidShortCode
    }

    if limit <= 0 || limit > 100 {
        limit = 10
    }

    clicks, err := s.analyticsStore.GetRecentClicks(ctx, shortCode, limit)
    if err != nil {
        return nil, fmt.Errorf("failed to get recent clicks: %w", err)
    }

    return clicks, nil
}