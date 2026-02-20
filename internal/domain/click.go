package domain

import (
    "time"
)

type ClickEvent struct {
    ID        int64     `json:"id" db:"id"`
    ShortCode string    `json:"short_code" db:"short_code"`
    UserAgent string    `json:"user_agent" db:"user_agent"`
    IP        string    `json:"ip" db:"ip"`
    Referer   string    `json:"referer" db:"referer"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type AnalyticsResponse struct {
    ShortCode    string                `json:"short_code"`
    OriginalURL  string                `json:"original_url"`
    TotalClicks  int64                 `json:"total_clicks"`
    DailyStats   map[string]int64       `json:"daily_stats"`
    MonthlyStats map[string]int64       `json:"monthly_stats"`
    Devices      map[string]int64       `json:"devices"`
    RecentClicks []ClickEvent           `json:"recent_clicks"`
}