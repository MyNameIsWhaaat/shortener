package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	pgxdriver "github.com/wb-go/wbf/dbpg/pgx-driver"
	wbflogger "github.com/wb-go/wbf/logger"

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

	appLogger, err := wbflogger.InitLogger(
		wbflogger.ZerologEngine,
		"shortener",
		"dev",
	)
	if err != nil {
		logger.Error("Failed to init wbf logger", "error", err)
		os.Exit(1)
	}

	logger.Info("Attempting to connect to database")
	pg, err := pgxdriver.New(cfg.PostgresDSN, appLogger)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pg.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pg.Ping(ctx); err != nil {
		logger.Error("Failed to ping database", "error", err)
		os.Exit(1)
	}

	pgStore := store.NewPostgresStoreFromPool(pg.Pool)
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

	sigCh := make(chan os.Signal, 1)
	errCh := make(chan error, 1)

	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	go func() {
		logger.Info("Server starting", "port", cfg.ServerPort)
		errCh <- server.Start()
	}()

	logger.Info("Server is running. Press Ctrl+C to stop.")

	select {
	case sig := <-sigCh:
		logger.Info("Shutdown signal received", "signal", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logger.Error("Failed to stop server gracefully", "error", err)
		}

	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Server stopped with error", "error", err)
			os.Exit(1)
		}
	}

	logger.Info("Application stopped")
}