package delete

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"url-shortener/internal/http-server/handlers/url/delete/mocks"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		mockError error
		respError string
	}{
		{
			name:  "Success",
			alias: "test-alias",
		},
		{
			name:      "Not Found",
			alias:     "non-existent",
			mockError: storage.ErrUrlNotFound,
			respError: "not found",
		},
		{
			name:      "Internal Error",
			alias:     "test-alias",
			mockError: errors.New("unexpected error"),
			respError: "internal error",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Init dependencies
			urlDeleterMock := mocks.NewURLDeleter(t)
			log := slogdiscard.NewDiscardLogger()

			// Setup mock expectations
			urlDeleterMock.On("DeleteURL", tc.alias).
				Return(tc.mockError).
				Once()

			// Init handler and router
			handler := New(log, urlDeleterMock)
			r := chi.NewRouter()
			r.Delete("/url/{alias}", handler)

			// Create request
			req, err := http.NewRequest(http.MethodDelete, "/url/"+tc.alias, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			// Execute request
			r.ServeHTTP(rr, req)

			// Assert results
			assert.Equal(t, http.StatusOK, rr.Code)

			if tc.respError == "" {
				// Assert success response
				assert.Contains(t, rr.Body.String(), `"status":"OK"`)
			} else {
				// Assert error response
				assert.Contains(t, rr.Body.String(), tc.respError)
			}

			// Verify mock calls
			urlDeleterMock.AssertExpectations(t)
		})
	}
}
