package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"syscall"

	"github.com/MyNameIsWhaaat/shortener/internal/api"
	"github.com/MyNameIsWhaaat/shortener/internal/config"
	handler "github.com/MyNameIsWhaaat/shortener/internal/httpapi"
	"github.com/MyNameIsWhaaat/shortener/internal/service"
	"github.com/MyNameIsWhaaat/shortener/internal/store"
)

func main() {
    cfg := config.Load()
    log.Printf("Config loaded: %+v", cfg)

    log.Println("Attempting to connect to database...")
    pgStore, err := store.NewPostgresStore(cfg.PostgresDSN)
    if err != nil {
        log.Fatalf("CRITICAL: Failed to connect to database: %v", err)
    }
    defer pgStore.Close()
    log.Println("Database connected successfully")

    shortenerService := service.NewShortenerService(
        pgStore,
        cfg.BaseURL,
        cfg.ShortCodeLength,
        pgStore,
    )
    
    analyticsService := service.NewAnalyticsService(pgStore)

    h := handler.NewHandler(shortenerService, analyticsService)

    server := api.NewServer(cfg, h)

    done := make(chan os.Signal, 1)
    signal.Notify(done, os.Interrupt, syscall.SIGTERM)

    go func() {
        log.Printf("Server starting on :%s", cfg.ServerPort)
        if err := server.Start(); err != nil {
            log.Printf("Server error: %v", err)
            if err != http.ErrServerClosed {
                log.Fatalf("Failed to start server: %v", err)
            }
        }
    }()

    log.Println("Server is running. Press Ctrl+C to stop.")
    <-done
    log.Println("Shutting down server...")
}