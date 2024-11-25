package userurls

import (
	"encoding/json"
	"github.com/orekhovskiy/shrtn/internal/entity"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/orekhovskiy/shrtn/internal/handler/http/api/mocks"
	"github.com/orekhovskiy/shrtn/internal/logger"
)

func TestGetUserURLs(t *testing.T) {
	mockLogger := &logger.NoopLogger{}
	mockAuthService := new(mocks.MockAuthService)
	mockURLService := new(mocks.MockURLService)

	handler := Handler{
		logger:      mockLogger,
		authService: mockAuthService,
		urlService:  mockURLService,
	}

	tests := []struct {
		name               string
		userID             string
		authServiceReturns bool
		urlServiceReturn   []entity.URLRecord
		urlServiceError    error
		expectedStatus     int
		expectedBody       string
	}{
		{
			name:               "Unauthorized when user ID is missing",
			authServiceReturns: false,
			expectedStatus:     http.StatusUnauthorized,
			expectedBody:       "Unauthorized",
		},
		{
			name:               "Internal Server Error from URL service",
			userID:             "testuser",
			authServiceReturns: true,
			urlServiceError:    assert.AnError,
			expectedStatus:     http.StatusInternalServerError,
			expectedBody:       "Internal Server Error",
		},
		{
			name:               "No URLs for user",
			userID:             "testuser",
			authServiceReturns: true,
			urlServiceReturn:   []entity.URLRecord{},
			expectedStatus:     http.StatusNoContent,
		},
		{
			name:               "Successful URL retrieval",
			userID:             "testuser",
			authServiceReturns: true,
			urlServiceReturn: []entity.URLRecord{
				{ShortURL: "http://short.ly/abc123", OriginalURL: "http://example.com"},
				{ShortURL: "http://short.ly/xyz789", OriginalURL: "http://another.com"},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"short_url":"http://short.ly/abc123","original_url":"http://example.com"},{"short_url":"http://short.ly/xyz789","original_url":"http://another.com"}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock calls
			mockAuthService.ExpectedCalls = nil
			mockURLService.ExpectedCalls = nil

			// Mock Auth Service
			mockAuthService.On("GetUserIDFromContext", mock.Anything).Return(tt.userID, tt.authServiceReturns)

			// Mock URL Service
			if tt.authServiceReturns {
				mockURLService.On("GetUserURLs", tt.userID).Return(tt.urlServiceReturn, tt.urlServiceError)
			}

			req := httptest.NewRequest(http.MethodGet, "/user/urls", nil)
			rec := httptest.NewRecorder()

			handler.getUserURLs(rec, req)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, rec.Code)

			// Assert response body if applicable
			if tt.expectedBody != "" {
				if rec.Header().Get("Content-Type") == "application/json" {
					var expected, actual []map[string]string
					if err := json.Unmarshal([]byte(tt.expectedBody), &expected); err != nil {
						t.Fatalf("Failed to unmarshal expected body: %v", err)
					}
					if err := json.Unmarshal(rec.Body.Bytes(), &actual); err != nil {
						t.Fatalf("Failed to unmarshal actual body: %v", err)
					}
					assert.Equal(t, expected, actual)
				} else {
					responseBody := strings.TrimSpace(rec.Body.String())
					assert.Equal(t, tt.expectedBody, responseBody)
				}
			}
		})
	}
}
