package shorten

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/orekhovskiy/shrtn/config"
	"github.com/orekhovskiy/shrtn/internal/handler/http/api/mocks"
	"github.com/orekhovskiy/shrtn/internal/logger"
)

func TestCreateShortUrl(t *testing.T) {
	mockLogger := &logger.NoopLogger{}
	mockService := new(mocks.MockURLService)
	opts := config.Config{BaseURL: "http://localhost:8080"}
	handler := Handler{logger: mockLogger, opts: opts, urlService: mockService}

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		mockServiceID  string
		contentType    string
	}{
		{
			name:           "Valid JSON request",
			requestBody:    `{"url":"https://example.com"}`,
			expectedStatus: http.StatusCreated,
			mockServiceID:  "abc123",
			contentType:    "application/json",
		},
		{
			name:           "Invalid JSON request",
			requestBody:    `{"invalid_json"`,
			expectedStatus: http.StatusBadRequest,
			contentType:    "application/json",
		},
		{
			name:           "Invalid URL",
			requestBody:    `{"url":"invalid_url"}`,
			expectedStatus: http.StatusBadRequest,
			contentType:    "application/json",
		},
		{
			name:           "Missing Content-Type",
			requestBody:    `{"url":"https://example.com"}`,
			expectedStatus: http.StatusBadRequest,
			contentType:    "",
		},
		{
			name:           "Wrong Content-Type",
			requestBody:    `{"url":"https://example.com"}`,
			expectedStatus: http.StatusBadRequest,
			contentType:    "text/plain",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock behavior for valid URL request
			if tt.expectedStatus == http.StatusCreated {
				mockService.On("Save", "https://example.com").Return(tt.mockServiceID)
			}

			// Create request
			req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer([]byte(tt.requestBody)))
			req.Header.Set("Content-Type", tt.contentType)

			// Record the response
			rec := httptest.NewRecorder()

			// Call the handler
			handler.CreateShortURL(rec, req)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rec.Code)

			// Check the response for the valid case
			if tt.expectedStatus == http.StatusCreated {
				var response ShortenResponse
				err := json.NewDecoder(rec.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, "http://localhost:8080/"+tt.mockServiceID, response.Result)
			}
		})
	}
}
