package handler

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    shortCode := vars["short_code"]
    
    log.Printf("Redirect request for short_code: %s", shortCode)

    url, err := h.shortenerService.GetOriginalURL(r.Context(), shortCode)
    if err != nil {
        log.Printf("URL not found: %v", err)
        http.Error(w, "URL not found", http.StatusNotFound)
        return
    }

    remoteAddr := r.RemoteAddr
    ip := remoteAddr
    if host, _, err := net.SplitHostPort(remoteAddr); err == nil {
        ip = host
    }
    
    log.Printf("RemoteAddr: %s -> IP: %s", remoteAddr, ip)

    go func() {
        ctx := context.Background()
        if err := h.shortenerService.TrackClick(ctx, shortCode, 
            r.UserAgent(), ip, r.Referer()); err != nil{
            log.Printf("Failed to track click: %v", err)
        } else {
            log.Printf("Click tracked successfully for %s", shortCode)
        }
    }()

    log.Printf("Redirecting to %s", url.OriginalURL)
    http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}
