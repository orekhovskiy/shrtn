package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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
		method         string
		contentType    string
		body           string
		mockSaveReturn string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Successful Short URL Creation",
			method:         http.MethodPost,
			contentType:    "text/plain",
			body:           "http://example.com",
			expectedStatus: http.StatusCreated,
			mockSaveReturn: "12345",
			expectedBody:   "http://localhost:8080/12345",
		},
		{
			name:           "Invalid URL in Body",
			method:         http.MethodPost,
			contentType:    "text/plain",
			body:           "invalid-url",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Bad Request\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set expectations for mock URL service
			if tt.mockSaveReturn != "" {
				mockURLService.On("Save", tt.body, mock.Anything).Return(tt.mockSaveReturn, nil)
			}

			mockAuthService.On("GetUserIDFromContext", mock.Anything).Return("testuser", true)
			mockURLService.On("BuildURL", "12345").Return("http://localhost:8080/12345")

			req := httptest.NewRequest(tt.method, "/shorten", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", tt.contentType)

			rec := httptest.NewRecorder()

			assert.NotNil(t, handler.urlService, "URL service should not be nil")
			assert.NotNil(t, handler.authService, "Auth service should not be nil")

			handler.CreateShortURL(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			responseBody, _ := io.ReadAll(rec.Body)
			assert.Equal(t, tt.expectedBody, string(responseBody))
		})
	}
}
