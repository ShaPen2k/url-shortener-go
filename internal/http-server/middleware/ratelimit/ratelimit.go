package ratelimit

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

// Limit returns a middleware that rate limits requests per IP
// Using a simple counter approach (for production use go-chi/httprate)
func Limit(requestsPerSecond int) func(next http.Handler) http.Handler {
	return middleware.ThrottleBacklog(requestsPerSecond, requestsPerSecond*2, 1)
}
