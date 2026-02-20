package handler

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/MyNameIsWhaaat/shortener/internal/logger"
	"github.com/gorilla/mux"
)

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortCode := vars["short_code"]

	logger.Info("Redirect request", "short_code", shortCode)

	url, err := h.shortenerService.GetOriginalURL(r.Context(), shortCode)
	if err != nil {
		logger.Error("URL not found", "short_code", shortCode, "error", err)
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	remoteAddr := r.RemoteAddr
	ip := remoteAddr
	if host, _, err := net.SplitHostPort(remoteAddr); err == nil {
		ip = host
	}

	logger.Info("Request IP", "remote_addr", remoteAddr, "ip", ip)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := h.shortenerService.TrackClick(ctx, shortCode,
			r.UserAgent(), ip, r.Referer()); err != nil {
			logger.Error("Failed to track click", "short_code", shortCode, "error", err)
		} else {
			logger.Info("Click tracked successfully", "short_code", shortCode)
		}
	}()

	logger.Info("Redirecting", "short_code", shortCode, "url", url.OriginalURL)
	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}
