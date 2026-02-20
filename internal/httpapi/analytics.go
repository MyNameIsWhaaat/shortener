package handler

import (
	"log"
	"net/http"

	"strconv"

	"github.com/MyNameIsWhaaat/shortener/internal/service"
	"github.com/gorilla/mux"
)

func (h *Handler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortCode := vars["short_code"]

	if shortCode == "" {
		h.respondError(w, "Short code is required", http.StatusBadRequest)
		return
	}

	analytics, err := h.analyticsService.GetAnalytics(r.Context(), shortCode)
	if err != nil {
		if service.IsNotFound(err) {
			h.respondError(w, "URL not found", http.StatusNotFound)
		} else {
			h.respondError(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.respond(w, analytics, http.StatusOK)
}

func (h *Handler) GetDailyStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortCode := vars["short_code"]

	log.Printf("GetDailyStats called for %s", shortCode)

	days := 30
	if daysStr := r.URL.Query().Get("days"); daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
			days = d
		}
	}

	stats, err := h.analyticsService.GetDailyStats(r.Context(), shortCode, days)
	if err != nil {
		if service.IsNotFound(err) {
			h.respondError(w, "URL not found", http.StatusNotFound)
		} else {
			h.respondError(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.respond(w, stats, http.StatusOK)
}

func (h *Handler) GetMonthlyStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortCode := vars["short_code"]

	log.Printf("GetMonthlyStats called for %s", shortCode)

	days := 30
	if daysStr := r.URL.Query().Get("days"); daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
			days = d
		}
	}

	stats, err := h.analyticsService.GetMonthlyStats(r.Context(), shortCode, days)
	if err != nil {
		if service.IsNotFound(err) {
			h.respondError(w, "URL not found", http.StatusNotFound)
		} else {
			h.respondError(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.respond(w, stats, http.StatusOK)
}

func (h *Handler) GetDeviceStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortCode := vars["short_code"]

	log.Printf("GetDeviceStats called for %s", shortCode)

	stats, err := h.analyticsService.GetDeviceStats(r.Context(), shortCode)
	if err != nil {
		if service.IsNotFound(err) {
			h.respondError(w, "URL not found", http.StatusNotFound)
		} else {
			h.respondError(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.respond(w, stats, http.StatusOK)
}
