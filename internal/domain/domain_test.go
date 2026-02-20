package domain

import (
	"testing"
	"time"
)

func TestURL(t *testing.T) {
	url := &URL{
		ID:          1,
		ShortCode:   "abc123",
		OriginalURL: "https://example.com",
		CreatedAt:   time.Now(),
		Clicks:      42,
	}

	if url.ShortCode != "abc123" {
		t.Errorf("expected ShortCode abc123, got %s", url.ShortCode)
	}

	if url.OriginalURL != "https://example.com" {
		t.Errorf("expected OriginalURL, got %s", url.OriginalURL)
	}

	if url.Clicks != 42 {
		t.Errorf("expected 42 clicks, got %d", url.Clicks)
	}
}

func TestURLWithCustomAlias(t *testing.T) {
	alias := "myalias"
	url := &URL{
		ID:          1,
		ShortCode:   "abc123",
		OriginalURL: "https://example.com",
		CustomAlias: &alias,
		CreatedAt:   time.Now(),
		Clicks:      0,
	}

	if url.CustomAlias == nil {
		t.Error("expected CustomAlias to be set")
	}

	if *url.CustomAlias != "myalias" {
		t.Errorf("expected CustomAlias myalias, got %s", *url.CustomAlias)
	}
}

func TestCreateURLRequest(t *testing.T) {
	req := &CreateURLRequest{
		URL: "https://example.com",
	}

	if req.URL != "https://example.com" {
		t.Errorf("expected URL https://example.com, got %s", req.URL)
	}

	if req.CustomAlias != nil {
		t.Error("expected CustomAlias to be nil")
	}
}

func TestCreateURLRequestWithAlias(t *testing.T) {
	alias := "myalias"
	req := &CreateURLRequest{
		URL:         "https://example.com",
		CustomAlias: &alias,
	}

	if req.CustomAlias == nil {
		t.Error("expected CustomAlias to be set")
	}

	if *req.CustomAlias != "myalias" {
		t.Errorf("expected CustomAlias myalias, got %s", *req.CustomAlias)
	}
}

func TestCreateURLResponse(t *testing.T) {
	resp := &CreateURLResponse{
		ShortCode:   "abc123",
		ShortURL:    "http://localhost:8080/s/abc123",
		OriginalURL: "https://example.com",
	}

	if resp.ShortCode != "abc123" {
		t.Errorf("expected ShortCode abc123, got %s", resp.ShortCode)
	}

	if resp.ShortURL != "http://localhost:8080/s/abc123" {
		t.Errorf("expected ShortURL, got %s", resp.ShortURL)
	}

	if resp.OriginalURL != "https://example.com" {
		t.Errorf("expected OriginalURL, got %s", resp.OriginalURL)
	}
}

func TestClickEvent(t *testing.T) {
	event := &ClickEvent{
		ID:        1,
		ShortCode: "abc123",
		UserAgent: "Mozilla/5.0",
		IP:        "192.168.1.1",
		Referer:   "https://google.com",
		CreatedAt: time.Now(),
	}

	if event.ShortCode != "abc123" {
		t.Errorf("expected ShortCode abc123, got %s", event.ShortCode)
	}

	if event.UserAgent != "Mozilla/5.0" {
		t.Errorf("expected UserAgent Mozilla/5.0, got %s", event.UserAgent)
	}

	if event.IP != "192.168.1.1" {
		t.Errorf("expected IP 192.168.1.1, got %s", event.IP)
	}
}

func TestAnalyticsResponse(t *testing.T) {
	resp := &AnalyticsResponse{
		ShortCode:    "abc123",
		OriginalURL:  "https://example.com",
		TotalClicks:  42,
		DailyStats:   map[string]int64{"2026-02-20": 10},
		MonthlyStats: map[string]int64{"2026-02": 42},
		Devices:      map[string]int64{"mobile": 25, "desktop": 17},
	}

	if resp.ShortCode != "abc123" {
		t.Errorf("expected ShortCode abc123, got %s", resp.ShortCode)
	}

	if resp.TotalClicks != 42 {
		t.Errorf("expected TotalClicks 42, got %d", resp.TotalClicks)
	}

	if dailyClicks, ok := resp.DailyStats["2026-02-20"]; !ok || dailyClicks != 10 {
		t.Error("expected daily stats entry")
	}

	if deviceClicks, ok := resp.Devices["mobile"]; !ok || deviceClicks != 25 {
		t.Error("expected device stats entry")
	}
}
