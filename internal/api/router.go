package api

import (
	"net/http"

	handler "github.com/MyNameIsWhaaat/shortener/internal/httpapi"
	"github.com/MyNameIsWhaaat/shortener/internal/ui"
	"github.com/gorilla/mux"
)

type Router struct {
	router  *mux.Router
	handler *handler.Handler
}

func NewRouter(h *handler.Handler) *Router {
	r := mux.NewRouter()

	r.Use(handler.LoggingMiddleware)
	r.Use(handler.RecoverMiddleware)

	router := &Router{
		router:  r,
		handler: h,
	}

	router.setupRoutes()

	return router
}

func (r *Router) setupRoutes() {
	api := r.router.PathPrefix("/api").Subrouter()
	{
		api.HandleFunc("/shorten", r.handler.Shorten).Methods("POST")
		api.HandleFunc("/urls", r.handler.GetAllURLs).Methods("GET")
		api.HandleFunc("/analytics/{short_code}", r.handler.GetAnalytics).Methods("GET")
		api.HandleFunc("/analytics/{short_code}/daily", r.handler.GetDailyStats).Methods("GET")
		api.HandleFunc("/analytics/{short_code}/monthly", r.handler.GetMonthlyStats).Methods("GET")
		api.HandleFunc("/analytics/{short_code}/devices", r.handler.GetDeviceStats).Methods("GET")
	}

	r.router.HandleFunc("/s/{short_code}", r.handler.Redirect).Methods("GET")
	r.router.HandleFunc("/health", r.handler.Health).Methods("GET")

	uiMux := http.NewServeMux()
	ui.Register(uiMux)
	r.router.PathPrefix("/ui").Handler(uiMux)
	r.router.PathPrefix("/").Handler(uiMux)
}

func (r *Router) GetHandler() http.Handler {
	return r.router
}
