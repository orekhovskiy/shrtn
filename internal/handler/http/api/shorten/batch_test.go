package shorten

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/orekhovskiy/shrtn/config"
	"github.com/orekhovskiy/shrtn/internal/entity"
	"github.com/orekhovskiy/shrtn/internal/handler/http/api/mocks"
	"github.com/orekhovskiy/shrtn/internal/logger"
	"github.com/stretchr/testify/mock"
)

func TestBatch(t *testing.T) {
	tests := []struct {
		name           string
		input          []entity.BatchRequest
		mockResponse   []entity.BatchResponse
		mockErr        error
		expectedCode   int
		expectedOutput []entity.BatchResponse
		userID         string
	}{
		{
			name: "valid batch request",
			input: []entity.BatchRequest{
				{CorrelationID: "1", OriginalURL: "http://example.com"},
				{CorrelationID: "2", OriginalURL: "http://example.org"},
			},
			mockResponse: []entity.BatchResponse{
				{CorrelationID: "1", ShortURL: "short1"},
				{CorrelationID: "2", ShortURL: "short2"},
			},
			mockErr:      nil,
			expectedCode: http.StatusCreated,
			expectedOutput: []entity.BatchResponse{
				{CorrelationID: "1", ShortURL: "short1"},
				{CorrelationID: "2", ShortURL: "short2"},
			},
			userID: "test-user-id",
		},
		{
			name:         "empty batch request",
			input:        []entity.BatchRequest{},
			expectedCode: http.StatusBadRequest,
			userID:       "test-user-id",
		},
		{
			name:         "invalid json",
			input:        nil, // Will simulate invalid JSON by passing an empty body
			expectedCode: http.StatusBadRequest,
			userID:       "test-user-id",
		},
		{
			name: "processing error",
			input: []entity.BatchRequest{
				{CorrelationID: "1", OriginalURL: "http://example.com"},
			},
			mockErr:      assert.AnError,
			expectedCode: http.StatusInternalServerError,
			userID:       "test-user-id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := &logger.NoopLogger{}
			mockService := new(mocks.MockURLService)
			mockAuthService := new(mocks.MockAuthService)
			opts := config.Config{BaseURL: "http://localhost:8080"}
			handler := Handler{
				logger:      mockLogger,
				opts:        opts,
				urlService:  mockService,
				authService: mockAuthService,
			}

			// Conditionally mock GetUserIDFromContext
			if tt.expectedCode != http.StatusBadRequest {
				mockAuthService.On("GetUserIDFromContext", mock.Anything).
					Return(tt.userID, true).Once()
			}

			// Mock ProcessBatch if applicable
			if tt.name == "valid batch request" {
				mockService.On("ProcessBatch", tt.input, tt.userID).Return(tt.mockResponse, nil)
			} else if tt.name == "processing error" {
				mockService.On("ProcessBatch", tt.input, tt.userID).Return(nil, tt.mockErr)
			}

			// Prepare request
			var requestBody []byte
			if tt.name == "invalid json" {
				requestBody = []byte("invalid json")
			} else {
				requestBody, _ = json.Marshal(tt.input)
			}

			req := httptest.NewRequest("POST", "/api/shorten/batch", bytes.NewBuffer(requestBody))
			rr := httptest.NewRecorder()

			// Call the Batch handler
			handler.Batch(rr, req)

			// Assert HTTP status code
			assert.Equal(t, tt.expectedCode, rr.Code)

			// Validate response for successful cases
			if tt.expectedCode == http.StatusCreated {
				var actualResponse []entity.BatchResponse
				err := json.Unmarshal(rr.Body.Bytes(), &actualResponse)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOutput, actualResponse)
			}

			// Verify expectations
			mockService.AssertExpectations(t)
			mockAuthService.AssertExpectations(t)
		})
	}
}
