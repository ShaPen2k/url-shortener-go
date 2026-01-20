package health

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"url-shortener/internal/lib/logger/handlers/slogdiscard"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthHandler(t *testing.T) {
	cases := []struct {
		name       string
		method     string
		path       string
		respStatus int
		respBody   string
	}{
		{
			name:       "Health OK",
			method:     http.MethodGet,
			path:       "/health",
			respStatus: http.StatusOK,
			respBody:   `"status":"ok"`,
		},
		{
			name:       "Health Ping Pong",
			method:     http.MethodGet,
			path:       "/health",
			respStatus: http.StatusOK,
			respBody:   `"ping":"pong"`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup dependencies
			log := slogdiscard.NewDiscardLogger()

			// Create handler
			handler := New(log)

			// Initialize chi router
			r := chi.NewRouter()
			r.Get("/health", handler)

			// Create request
			req, err := http.NewRequest(tc.method, tc.path, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			// Execute request
			r.ServeHTTP(rr, req)

			// Assert status code
			assert.Equal(t, tc.respStatus, rr.Code)

			// Assert response body contains expected fields
			assert.Contains(t, rr.Body.String(), tc.respBody)

			// Assert content type
			assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
		})
	}
}

func TestHealthHandlerJSON(t *testing.T) {
	log := slogdiscard.NewDiscardLogger()
	handler := New(log)

	r := chi.NewRouter()
	r.Get("/health", handler)

	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Assert full JSON response structure
	expectedJSON := `{"status":"ok","ping":"pong"}`
	assert.JSONEq(t, expectedJSON, rr.Body.String())
}
