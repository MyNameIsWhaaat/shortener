package store

import (
	"context"
	"fmt"
	"log"

	"github.com/MyNameIsWhaaat/shortener/internal/domain"
)

func (s *PostgresStore) SaveClickEvent(ctx context.Context, event *domain.ClickEvent) error {
    log.Printf("Saving click event for %s", event.ShortCode)
    
    query := `
        INSERT INTO click_events (short_code, user_agent, ip, referer, created_at)
        VALUES ($1, $2, $3, $4, $5)
    `
    
    _, err := s.db.ExecContext(ctx, query,
        event.ShortCode,
        event.UserAgent,
        event.IP,
        event.Referer,
        event.CreatedAt,
    )
    
    if err != nil {
        log.Printf("Failed to save click event: %v", err)
        return err
    }
    
    log.Printf("Click event saved for %s", event.ShortCode)
    return nil
}

func (s *PostgresStore) GetAnalytics(ctx context.Context, shortCode string) (*domain.AnalyticsResponse, error) {
    url, err := s.GetURLByShortCode(ctx, shortCode)
    if err != nil {
        return nil, err
    }

    response := &domain.AnalyticsResponse{
        ShortCode:    url.ShortCode,
        OriginalURL:  url.OriginalURL,
        TotalClicks:  url.Clicks,
    }

    response.DailyStats, _ = s.GetDailyStats(ctx, shortCode, 30)
    response.MonthlyStats, _ = s.GetMonthlyStats(ctx, shortCode, 12)
    response.Devices, _ = s.GetDeviceStats(ctx, shortCode)
    response.RecentClicks, _ = s.GetRecentClicks(ctx, shortCode, 10)

    return response, nil
}

func (s *PostgresStore) GetDailyStats(ctx context.Context, shortCode string, days int) (map[string]int64, error) {
    log.Printf("Store GetDailyStats for %s, days=%d", shortCode, days)

    query := `
        SELECT TO_CHAR(created_at, 'YYYY-MM-DD') as day, COUNT(*) 
        FROM click_events 
        WHERE short_code = $1 
        AND created_at >= NOW() - ($2 || ' days')::INTERVAL
        GROUP BY day 
        ORDER BY day DESC
    `
    
    rows, err := s.db.QueryContext(ctx, query, shortCode, days)
    if err != nil {
        return nil, fmt.Errorf("failed to get daily stats: %w", err)
    }
    defer rows.Close()

    stats := make(map[string]int64)
    for rows.Next() {
        var day string
        var count int64
        if err := rows.Scan(&day, &count); err != nil {
            return nil, fmt.Errorf("failed to scan daily stats: %w", err)
        }
        stats[day] = count
    }
    return stats, nil
}

func (s *PostgresStore) GetMonthlyStats(ctx context.Context, shortCode string, months int) (map[string]int64, error) {
    query := `
        SELECT TO_CHAR(created_at, 'YYYY-MM') as month, COUNT(*) 
        FROM click_events 
        WHERE short_code = $1 
        AND created_at >= NOW() - ($2 || ' months')::INTERVAL
        GROUP BY month 
        ORDER BY month DESC
    `
    
    rows, err := s.db.QueryContext(ctx, query, shortCode, months)
    if err != nil {
        return nil, fmt.Errorf("failed to get monthly stats: %w", err)
    }
    defer rows.Close()

    stats := make(map[string]int64)
    for rows.Next() {
        var month string
        var count int64
        if err := rows.Scan(&month, &count); err != nil {
            return nil, fmt.Errorf("failed to scan monthly stats: %w", err)
        }
        stats[month] = count
    }
    return stats, nil
}

func (s *PostgresStore) GetDeviceStats(ctx context.Context, shortCode string) (map[string]int64, error) {
    query := `
        SELECT 
            CASE 
                WHEN user_agent ILIKE '%mobile%' OR user_agent ILIKE '%android%' OR user_agent ILIKE '%iphone%' THEN 'Mobile'
                WHEN user_agent ILIKE '%tablet%' OR user_agent ILIKE '%ipad%' THEN 'Tablet'
                WHEN user_agent ILIKE '%bot%' OR user_agent ILIKE '%crawler%' THEN 'Bot'
                ELSE 'Desktop'
            END as device_type,
            COUNT(*) 
        FROM click_events 
        WHERE short_code = $1 
        GROUP BY device_type
    `
    
    rows, err := s.db.QueryContext(ctx, query, shortCode)
    if err != nil {
        return nil, fmt.Errorf("failed to get device stats: %w", err)
    }
    defer rows.Close()

    stats := make(map[string]int64)
    for rows.Next() {
        var device string
        var count int64
        if err := rows.Scan(&device, &count); err != nil {
            return nil, fmt.Errorf("failed to scan device stats: %w", err)
        }
        stats[device] = count
    }
    return stats, nil
}

func (s *PostgresStore) GetRecentClicks(ctx context.Context, shortCode string, limit int) ([]domain.ClickEvent, error) {
    query := `
        SELECT user_agent, ip, referer, created_at
        FROM click_events 
        WHERE short_code = $1 
        ORDER BY created_at DESC 
        LIMIT $2
    `
    
    rows, err := s.db.QueryContext(ctx, query, shortCode, limit)
    if err != nil {
        return nil, fmt.Errorf("failed to get recent clicks: %w", err)
    }
    defer rows.Close()

    var clicks []domain.ClickEvent
    for rows.Next() {
        var event domain.ClickEvent
        event.ShortCode = shortCode
        if err := rows.Scan(&event.UserAgent, &event.IP, &event.Referer, &event.CreatedAt); err != nil {
            return nil, fmt.Errorf("failed to scan recent click: %w", err)
        }
        clicks = append(clicks, event)
    }
    return clicks, nil
}