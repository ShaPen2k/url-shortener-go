package redirect

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"url-shortener/internal/http-server/handlers/redirect/mocks"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedirectHandler(t *testing.T) {
	cases := []struct {
		name       string
		alias      string
		url        string
		mockError  error
		respStatus int
	}{
		{
			name:       "Success",
			alias:      "test-alias",
			url:        "https://google.com",
			respStatus: http.StatusFound, // 302
		},
		{
			name:       "URL Not Found",
			alias:      "non-existent",
			mockError:  storage.ErrUrlNotFound,
			respStatus: http.StatusNotFound, // 200
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// initialize mock
			urlGetterMock := mocks.NewURLGetter(t)

			if tc.respStatus == http.StatusFound || tc.mockError != nil {
				urlGetterMock.On("GetURL", tc.alias).Return(tc.url, tc.mockError).Once()
			}

			// create fakeLogger
			log := slogdiscard.NewDiscardLogger()

			// create handler
			handler := New(log, urlGetterMock)

			// initialize chi router
			r := chi.NewRouter()
			r.Get("/{alias}", handler)

			// create request
			req, err := http.NewRequest(http.MethodGet, "/"+tc.alias, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			// fill alias context
			r.ServeHTTP(rr, req)

			// check status
			assert.Equal(t, tc.respStatus, rr.Code)

			if tc.respStatus == http.StatusFound {
				// check redirect (Location)
				assert.Equal(t, tc.url, rr.Header().Get("Location"))
			} else {
				// check error in JSON
				assert.Contains(t, rr.Body.String(), "not found")
			}
		})
	}
}
