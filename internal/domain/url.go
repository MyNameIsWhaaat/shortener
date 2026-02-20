package domain

import (
	"time"
)

type URL struct {
	ID          int64     `json:"id" db:"id"`
	ShortCode   string    `json:"short_code" db:"short_code"`
	OriginalURL string    `json:"original_url" db:"original_url"`
	CustomAlias *string   `json:"custom_alias,omitempty" db:"custom_alias"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	Clicks      int64     `json:"clicks" db:"clicks"`
}

type CreateURLRequest struct {
	URL         string  `json:"url"`
	CustomAlias *string `json:"custom_alias,omitempty"`
}

type CreateURLResponse struct {
	ShortCode   string `json:"short_code"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
