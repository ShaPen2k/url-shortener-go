package health

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

type HealthResponse struct {
	Status string `json:"status"`
	Ping   string `json:"ping"`
}

// New returns a handler that checks service health
func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		render.JSON(w, r, HealthResponse{
			Status: "ok",
			Ping:   "pong",
		})
	}
}
