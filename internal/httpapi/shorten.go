package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/MyNameIsWhaaat/shortener/internal/domain"
	"github.com/MyNameIsWhaaat/shortener/internal/service"
)

type shortenRequest struct {
    URL         string  `json:"url"`
    CustomAlias *string `json:"custom_alias,omitempty"`
}

type shortenResponse struct {
    ShortCode   string `json:"short_code"`
    ShortURL    string `json:"short_url"`
    OriginalURL string `json:"original_url"`
}

func (h *Handler) GetAllURLs(w http.ResponseWriter, r *http.Request) {
    // Получаем limit из query параметра
    limit := 50
    if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
        if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
            limit = l
        }
    }
    
    urls, err := h.shortenerService.GetAllURLs(r.Context(), limit)
    if err != nil {
        h.respondError(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    h.respond(w, urls, http.StatusOK)
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
    var req shortenRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.respondError(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    createReq := &domain.CreateURLRequest{
        URL:         req.URL,
        CustomAlias: req.CustomAlias,
    }

    resp, err := h.shortenerService.CreateShortURL(r.Context(), createReq)
    if err != nil {
        switch {
        case service.IsAlreadyExists(err):
            h.respondError(w, "Short code already exists", http.StatusConflict)
        case err == service.ErrInvalidURL:
            h.respondError(w, "Invalid URL", http.StatusBadRequest)
        case err == service.ErrEmptyURL:
            h.respondError(w, "URL cannot be empty", http.StatusBadRequest)
        case err == service.ErrInvalidShortCode:
            h.respondError(w, "Invalid custom short code", http.StatusBadRequest)
        default:
            h.respondError(w, "Internal server error", http.StatusInternalServerError)
        }
        return
    }

    h.respond(w, shortenResponse{
        ShortCode:   resp.ShortCode,
        ShortURL:    resp.ShortURL,
        OriginalURL: resp.OriginalURL,
    }, http.StatusCreated)
}