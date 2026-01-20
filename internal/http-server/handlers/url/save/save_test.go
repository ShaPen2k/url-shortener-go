package save

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"url-shortener/internal/lib/logger/handlers/slogdiscard"
	"url-shortener/internal/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// urlSaverMock is a mock implementation of the URLSaver interface
type urlSaverMock struct {
	mock.Mock
}

func (m *urlSaverMock) SaveURL(urlToSave string, alias string) (int64, error) {
	args := m.Called(urlToSave, alias)
	return args.Get(0).(int64), args.Error(1)
}

func TestSaveHandler(t *testing.T) {
	// Define test cases for table-driven testing
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test-alias",
			url:   "https://google.com",
		},
		{
			name:  "Empty Alias",
			alias: "",
			url:   "https://yandex.ru",
		},
		{
			name:      "Empty URL",
			url:       "",
			respError: "field URL is a required field",
		},
		{
			name:      "Invalid URL",
			url:       "not-a-url",
			respError: "field URL is not a valid URL",
		},
		{
			name:      "Save Error",
			alias:     "fail-alias",
			url:       "https://fail.com",
			mockError: errors.New("unexpected error"),
			respError: "failed to add url",
		},
		{
			name:      "URL Already Exists",
			alias:     "exists",
			url:       "https://exists.com",
			mockError: storage.ErrUrlExists,
			respError: "url already exists",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup dependencies
			urlSaverMock := new(urlSaverMock)
			log := slogdiscard.NewDiscardLogger()

			// Setup mock expectations
			if tc.respError == "" || tc.mockError != nil {
				urlSaverMock.On("SaveURL", tc.url, mock.AnythingOfType("string")).
					Return(int64(1), tc.mockError).
					Once()
			}

			// Init handler
			handler := New(log, urlSaverMock)

			// Prepare request body
			input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tc.url, tc.alias)
			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			// Execute request
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			// Assert HTTP status
			assert.Equal(t, http.StatusOK, rr.Code)

			// Parse and assert response body
			var resp Response
			err = json.Unmarshal(rr.Body.Bytes(), &resp)
			require.NoError(t, err)

			if tc.respError == "" {
				// Assert successful response
				assert.Equal(t, "OK", resp.Status)
				assert.Empty(t, resp.Error)
				assert.NotEmpty(t, resp.Alias)
				if tc.alias != "" {
					assert.Equal(t, tc.alias, resp.Alias)
				}
			} else {
				// Assert error response
				assert.Contains(t, resp.Error, tc.respError)
			}
		})
	}
}
