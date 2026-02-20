package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/MyNameIsWhaaat/shortener/internal/domain"
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
	if err := json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"time":   time.Now().String(),
	}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) GetPopularURLs(w http.ResponseWriter, r *http.Request) {
	limit := 10
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	urls, err := h.shortenerService.GetPopularURLs(r.Context(), limit)
	if err != nil {
		h.respondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if urls == nil {
		urls = []*domain.URL{}
	}

	h.respond(w, urls, http.StatusOK)
}
