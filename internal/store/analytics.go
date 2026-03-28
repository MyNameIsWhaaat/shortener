package store

import (
	"context"
	"fmt"

	"github.com/MyNameIsWhaaat/shortener/internal/domain"
	"github.com/MyNameIsWhaaat/shortener/internal/logger"
)

func (s *PostgresStore) SaveClickEvent(ctx context.Context, event *domain.ClickEvent) error {
	logger.Info("Saving click event", "short_code", event.ShortCode)

	query := `
        INSERT INTO click_events (short_code, user_agent, ip, referer, created_at)
        VALUES ($1, $2, $3, $4, $5)
    `

	_, err := s.db.Exec(ctx, query,
		event.ShortCode,
		event.UserAgent,
		event.IP,
		event.Referer,
		event.CreatedAt,
	)

	if err != nil {
		logger.Error("Failed to save click event", "error", err)
		return err
	}

	logger.Info("Click event saved", "short_code", event.ShortCode)
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
		CreatedAt:    url.CreatedAt,
		TotalClicks:  url.Clicks,
		DailyStats:   map[string]int64{},
		MonthlyStats: map[string]int64{},
		Devices:      map[string]int64{},
		RecentClicks: []domain.ClickEvent{},
	}

	if response.DailyStats, err = s.GetDailyStats(ctx, shortCode, 30); err != nil {
		logger.Error("GetDailyStats failed", "short_code", shortCode, "error", err)
		response.DailyStats = map[string]int64{}
	}

	if response.MonthlyStats, err = s.GetMonthlyStats(ctx, shortCode, 12); err != nil {
		logger.Error("GetMonthlyStats failed", "short_code", shortCode, "error", err)
		response.MonthlyStats = map[string]int64{}
	}

	if response.Devices, err = s.GetDeviceStats(ctx, shortCode); err != nil {
		logger.Error("GetDeviceStats failed", "short_code", shortCode, "error", err)
		response.Devices = map[string]int64{}
	}

	if response.RecentClicks, err = s.GetRecentClicks(ctx, shortCode, 10); err != nil {
		logger.Error("GetRecentClicks failed", "short_code", shortCode, "error", err)
		response.RecentClicks = []domain.ClickEvent{}
	}

	return response, nil
}

func (s *PostgresStore) GetDailyStats(ctx context.Context, shortCode string, days int) (map[string]int64, error) {
	query := `
        SELECT TO_CHAR(created_at::timestamp, 'YYYY-MM-DD') AS day, COUNT(*)::bigint
        FROM click_events
        WHERE short_code = $1
          AND created_at >= NOW() - INTERVAL '30 days'
        GROUP BY TO_CHAR(created_at::timestamp, 'YYYY-MM-DD')
        ORDER BY day DESC
    `

	rows, err := s.db.Query(ctx, query, shortCode)
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

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("daily stats rows error: %w", err)
	}

	return stats, nil
}

func (s *PostgresStore) GetMonthlyStats(ctx context.Context, shortCode string, months int) (map[string]int64, error) {
	query := `
        SELECT TO_CHAR(created_at::timestamp, 'YYYY-MM') AS month, COUNT(*)::bigint
        FROM click_events
        WHERE short_code = $1
          AND created_at >= NOW() - INTERVAL '12 months'
        GROUP BY TO_CHAR(created_at::timestamp, 'YYYY-MM')
        ORDER BY month DESC
    `

	rows, err := s.db.Query(ctx, query, shortCode)
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

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("monthly stats rows error: %w", err)
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

	rows, err := s.db.Query(ctx, query, shortCode)
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

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("device stats rows error: %w", err)
	}

	return stats, nil
}

func (s *PostgresStore) GetRecentClicks(ctx context.Context, shortCode string, limit int) ([]domain.ClickEvent, error) {
	query := `
        SELECT user_agent, ip::text, COALESCE(referer, ''), created_at
        FROM click_events
        WHERE short_code = $1
        ORDER BY created_at DESC
        LIMIT $2
    `

	rows, err := s.db.Query(ctx, query, shortCode, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent clicks: %w", err)
	}
	defer rows.Close()

	clicks := make([]domain.ClickEvent, 0)
	for rows.Next() {
		var event domain.ClickEvent
		event.ShortCode = shortCode

		if err := rows.Scan(
			&event.UserAgent,
			&event.IP,
			&event.Referer,
			&event.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan recent click: %w", err)
		}

		clicks = append(clicks, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("recent clicks rows error: %w", err)
	}

	return clicks, nil
}