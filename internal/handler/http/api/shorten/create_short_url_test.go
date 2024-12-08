package shorten

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/orekhovskiy/shrtn/config"
	"github.com/orekhovskiy/shrtn/internal/handler/http/api/mocks"
	"github.com/orekhovskiy/shrtn/internal/logger"
)

func TestCreateShortUrl(t *testing.T) {
	mockLogger := &logger.NoopLogger{}
	mockURLService := new(mocks.MockURLService)
	mockAuthService := new(mocks.MockAuthService)

	opts := config.Config{BaseURL: "http://localhost:8080"}
	handler := Handler{
		logger:      mockLogger,
		opts:        opts,
		urlService:  mockURLService,
		authService: mockAuthService,
	}

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		mockServiceID  string
		contentType    string
		userID         string
	}{
		{
			name:           "Valid JSON request",
			requestBody:    `{"url":"https://example.com"}`,
			expectedStatus: http.StatusCreated,
			mockServiceID:  "abc123",
			contentType:    "application/json",
			userID:         "test-user-id",
		},
		{
			name:           "Invalid JSON request",
			requestBody:    `{"invalid_json"`,
			expectedStatus: http.StatusBadRequest,
			contentType:    "application/json",
			userID:         "test-user-id",
		},
		{
			name:           "Invalid URL",
			requestBody:    `{"url":"invalid_url"}`,
			expectedStatus: http.StatusBadRequest,
			contentType:    "application/json",
			userID:         "test-user-id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock behavior for user ID extraction
			mockAuthService.On("GetUserIDFromContext", mock.Anything).
				Return(tt.userID, true).Once()

			// Mock behavior for URL service
			if tt.expectedStatus == http.StatusCreated {
				mockURLService.On("Save", "https://example.com", tt.userID).
					Return(tt.mockServiceID, nil).Once()

				// Ensure BuildURL returns the correct value with no error
				mockURLService.On("BuildURL", tt.mockServiceID).
					Return(opts.BaseURL+"/"+tt.mockServiceID, nil).Once() // Correct return value with no error
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

			// Validate the response for the valid case
			if tt.expectedStatus == http.StatusCreated {
				var response ShortenResponse
				err := json.NewDecoder(rec.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, "http://localhost:8080/"+tt.mockServiceID, response.Result)
			}
		})
	}
}
