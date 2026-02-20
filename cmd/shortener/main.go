package main

import (
	"net/http"
	"os"
	"os/signal"

	"syscall"

	"github.com/MyNameIsWhaaat/shortener/internal/api"
	"github.com/MyNameIsWhaaat/shortener/internal/cache"
	"github.com/MyNameIsWhaaat/shortener/internal/config"
	handler "github.com/MyNameIsWhaaat/shortener/internal/httpapi"
	"github.com/MyNameIsWhaaat/shortener/internal/logger"
	"github.com/MyNameIsWhaaat/shortener/internal/service"
	"github.com/MyNameIsWhaaat/shortener/internal/store"
)

func main() {
	logger.Init()

	cfg := config.Load()
	logger.Info("Config loaded", "port", cfg.ServerPort, "base_url", cfg.BaseURL)

	logger.Info("Attempting to connect to database")
	pgStore, err := store.NewPostgresStore(cfg.PostgresDSN)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pgStore.Close()
	logger.Info("Database connected successfully")

	var cacheClient cache.Cache
	redisCache, err := cache.NewRedisCache(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB, cfg.CacheTTL)
	if err != nil {
		logger.Warn("Failed to initialize Redis cache, continuing without cache", "error", err)
		cacheClient = &cache.NoOpCache{}
	} else {
		defer redisCache.Close()
		cacheClient = redisCache
		logger.Info("Redis cache initialized successfully")
	}

	shortenerService := service.NewShortenerService(
		pgStore,
		cfg.BaseURL,
		cfg.ShortCodeLength,
		pgStore,
		cacheClient,
	)

	analyticsService := service.NewAnalyticsService(pgStore)

	h := handler.NewHandler(shortenerService, analyticsService)

	server := api.NewServer(cfg, h)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Info("Server starting", "port", cfg.ServerPort)
		if err := server.Start(); err != nil {
			logger.Error("Server error", "error", err)
			if err != http.ErrServerClosed {
				os.Exit(1)
			}
		}
	}()

	logger.Info("Server is running. Press Ctrl+C to stop.")
	<-done
	logger.Info("Shutting down server")
}

