package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"regexp"

	"time"

	"github.com/MyNameIsWhaaat/shortener/internal/cache"
	"github.com/MyNameIsWhaaat/shortener/internal/domain"
	"github.com/MyNameIsWhaaat/shortener/internal/store"
)

type shortenerService struct {
    urlStore store.URLStore
    baseURL  string
    codeLen  int
    analyticsStore store.AnalyticsStore
    cache cache.Cache
}

func NewShortenerService(urlStore store.URLStore, baseURL string, codeLen int, analyticsStore store.AnalyticsStore, cacheClient cache.Cache) ShortenerService {
    return &shortenerService{
        urlStore: urlStore,
        baseURL:  baseURL,
        codeLen:  codeLen,
        analyticsStore: analyticsStore,
        cache: cacheClient,
    }
}

func (s *shortenerService) GetAllURLs(ctx context.Context, limit int) ([]*domain.URL, error) {
    if limit <= 0 || limit > 100 {
        limit = 50
    }
    
    urls, err := s.urlStore.GetAllURLs(ctx, limit)
    if err != nil {
        return nil, fmt.Errorf("failed to get all urls: %w", err)
    }
    
    return urls, nil
}

func (s *shortenerService) GetPopularURLs(ctx context.Context, limit int) ([]*domain.URL, error) {
    if s.cache == nil {
        return nil, fmt.Errorf("cache not available")
    }

    if limit <= 0 || limit > 100 {
        limit = 10
    }

    urls, err := s.cache.GetPopular(ctx, limit)
    if err != nil {
        return nil, fmt.Errorf("failed to get popular urls: %w", err)
    }

    return urls, nil
}

func (s *shortenerService) CreateShortURL(ctx context.Context, req *domain.CreateURLRequest) (*domain.CreateURLResponse, error) {
    if err := s.validateURL(req.URL); err != nil {
        return nil, fmt.Errorf("invalid URL: %w", err)
    }

    shortCode := req.CustomAlias
    if shortCode == nil {
        code, err := s.generateShortCode()
        if err != nil {
            return nil, fmt.Errorf("failed to generate short code: %w", err)
        }
        shortCode = &code
    } else {
        if err := s.validateShortCode(*shortCode); err != nil {
            return nil, fmt.Errorf("invalid custom short code: %w", err)
        }
    }

    exists, err := s.urlStore.CheckShortCodeExists(ctx, *shortCode)
    if err != nil {
        return nil, fmt.Errorf("failed to check code existence: %w", err)
    }
    if exists {
        return nil, ErrShortCodeExists
    }

    url := &domain.URL{
        ShortCode:   *shortCode,
        OriginalURL: req.URL,
        CustomAlias: req.CustomAlias,
        CreatedAt:   time.Now(),
        Clicks:      0,
    }

    if err := s.urlStore.CreateURL(ctx, url); err != nil {
        return nil, fmt.Errorf("failed to create url in store: %w", err)
    }

    // Cache the newly created URL
    if s.cache != nil {
        if err := s.cache.Set(ctx, url.ShortCode, url); err != nil {
            log.Printf("Failed to cache newly created URL: %v", err)
        }
    }

    return &domain.CreateURLResponse{
        ShortCode:   url.ShortCode,
        ShortURL:    fmt.Sprintf("%s/s/%s", s.baseURL, url.ShortCode),
        OriginalURL: url.OriginalURL,
    }, nil
}

func (s *shortenerService) GetOriginalURL(ctx context.Context, shortCode string) (*domain.URL, error) {
    if shortCode == "" {
        return nil, ErrInvalidShortCode
    }
    
    // Try cache first
    if s.cache != nil {
        if cachedURL, err := s.cache.Get(ctx, shortCode); err == nil && cachedURL != nil {
            return cachedURL, nil
        }
    }

    // Fallback to database
    url, err := s.urlStore.GetURLByShortCode(ctx, shortCode)
    if err != nil {
        return nil, fmt.Errorf("failed to get url from store: %w", err)
    }

    // Cache the result
    if s.cache != nil {
        if err := s.cache.Set(ctx, shortCode, url); err != nil {
            log.Printf("Failed to cache URL: %v", err)
        }
    }

    return url, nil
}

func (s *shortenerService) TrackClick(ctx context.Context, shortCode, userAgent, ip, referer string) error {
    log.Printf("TrackClick for %s", shortCode)

    if err := s.urlStore.IncrementClicks(ctx, shortCode); err != nil {
        return err
    }
    
    event := &domain.ClickEvent{
        ShortCode: shortCode,
        UserAgent: userAgent,
        IP:        ip,
        Referer:   referer,
        CreatedAt: time.Now(),
    }
    
    if s.analyticsStore == nil {
        log.Printf("analyticsStore is nil!")
        return nil
    }
    
    return s.analyticsStore.SaveClickEvent(ctx, event)
}

func (s *shortenerService) generateShortCode() (string, error) {
    b := make([]byte, s.codeLen)
    if _, err := rand.Read(b); err != nil {
        return "", err
    }
    code := base64.URLEncoding.EncodeToString(b)
    return code[:s.codeLen], nil
}

func (s *shortenerService) validateURL(rawURL string) error {
    if rawURL == "" {
        return ErrEmptyURL
    }

    parsed, err := url.Parse(rawURL)
    if err != nil {
        return fmt.Errorf("failed to parse URL: %w", err)
    }

    if parsed.Scheme != "http" && parsed.Scheme != "https" {
        return fmt.Errorf("unsupported URL scheme: %s", parsed.Scheme)
    }

    if parsed.Host == "" {
        return fmt.Errorf("missing host in URL")
    }

    return nil
}

func (s *shortenerService) validateShortCode(code string) error {
    if code == "" {
        return ErrInvalidShortCode
    }

    if len(code) > 50 {
        return ErrShortCodeTooLong
    }

    matched, err := regexp.MatchString("^[a-zA-Z0-9_-]+$", code)
    if err != nil || !matched {
        return ErrInvalidShortCode
    }

    return nil
}