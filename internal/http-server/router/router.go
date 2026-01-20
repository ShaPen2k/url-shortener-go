// Package router provides HTTP routing configuration for the application.
package router

import (
	"log/slog"
	"url-shortener/internal/config"
	"url-shortener/internal/http-server/handlers/redirect"
	"url-shortener/internal/http-server/handlers/url/delete"
	"url-shortener/internal/http-server/handlers/url/save"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Storage is a composite interface that groups all storage-related requirements
// needed for the HTTP handlers.
type Storage interface {
	save.URLSaver
	redirect.URLGetter
	delete.URLDeleter
}

// Setup initializes the chi router with global middleware and application routes.
func Setup(log *slog.Logger, cfg config.HTTPServer, storage Storage) *chi.Mux {
	r := chi.NewRouter()

	// Apply standard middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	// Protected routes (require Basic Auth)
	r.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.User: cfg.Password,
		}))

		r.Post("/", save.New(log, storage))
		r.Delete("/{alias}", delete.New(log, storage))
	})

	// Public route for URL redirection
	r.Get("/{alias}", redirect.New(log, storage))

	return r
}
