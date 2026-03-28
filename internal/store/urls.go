package store

import (
	"context"
	"fmt"

	"github.com/MyNameIsWhaaat/shortener/internal/domain"
	"github.com/jackc/pgx/v5"
)

func (s *PostgresStore) CreateURL(ctx context.Context, url *domain.URL) error {
	query := `
        INSERT INTO urls (short_code, original_url, custom_alias, created_at, clicks)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `

	err := s.db.QueryRow(
		ctx,
		query,
		url.ShortCode,
		url.OriginalURL,
		url.CustomAlias,
		url.CreatedAt,
		url.Clicks,
	).Scan(&url.ID)

	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrShortCodeExists
		}
		return fmt.Errorf("failed to create url: %w", err)
	}

	return nil
}

func (s *PostgresStore) GetURLByShortCode(ctx context.Context, shortCode string) (*domain.URL, error) {
	query := `
        SELECT id, short_code, original_url, custom_alias, created_at, clicks
        FROM urls
        WHERE short_code = $1
    `

	var url domain.URL
	err := s.db.QueryRow(ctx, query, shortCode).Scan(
		&url.ID,
		&url.ShortCode,
		&url.OriginalURL,
		&url.CustomAlias,
		&url.CreatedAt,
		&url.Clicks,
	)

	if err == pgx.ErrNoRows {
		return nil, domain.ErrURLNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get url: %w", err)
	}

	return &url, nil
}

func (s *PostgresStore) IncrementClicks(ctx context.Context, shortCode string) error {
	query := `UPDATE urls SET clicks = clicks + 1 WHERE short_code = $1`

	tag, err := s.db.Exec(ctx, query, shortCode)
	if err != nil {
		return fmt.Errorf("failed to increment clicks: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrURLNotFound
	}

	return nil
}

func (s *PostgresStore) CheckShortCodeExists(ctx context.Context, shortCode string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM urls WHERE short_code = $1)`

	var exists bool
	err := s.db.QueryRow(ctx, query, shortCode).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check short code: %w", err)
	}

	return exists, nil
}

func (s *PostgresStore) GetAllURLs(ctx context.Context, limit int) ([]*domain.URL, error) {
	query := `
        SELECT id, short_code, original_url, custom_alias, created_at, clicks
        FROM urls
        ORDER BY created_at DESC
        LIMIT $1
    `

	rows, err := s.db.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get all urls: %w", err)
	}
	defer rows.Close()

	var urls []*domain.URL
	for rows.Next() {
		var url domain.URL
		if err := rows.Scan(
			&url.ID,
			&url.ShortCode,
			&url.OriginalURL,
			&url.CustomAlias,
			&url.CreatedAt,
			&url.Clicks,
		); err != nil {
			return nil, fmt.Errorf("failed to scan url: %w", err)
		}
		urls = append(urls, &url)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return urls, nil
}