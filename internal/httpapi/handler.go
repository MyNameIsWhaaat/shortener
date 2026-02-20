package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/MyNameIsWhaaat/shortener/internal/service"
)

type Handler struct {
    shortenerService service.ShortenerService
    analyticsService service.AnalyticsService
}

func NewHandler(
    shortenerService service.ShortenerService,
    analyticsService service.AnalyticsService,
) *Handler {
    return &Handler{
        shortenerService: shortenerService,
        analyticsService: analyticsService,
    }
}

func (h *Handler) respond(w http.ResponseWriter, data interface{}, status int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    if err := json.NewEncoder(w).Encode(data); err != nil {
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
    }
}

func (h *Handler) respondError(w http.ResponseWriter, message string, status int) {
    h.respond(w, map[string]string{"error": message}, status)
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "status": "ok",
        "time":   time.Now().String(),
    })
}