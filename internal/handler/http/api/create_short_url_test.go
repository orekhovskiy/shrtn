package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/orekhovskiy/shrtn/internal/handler/http/api/mocks"

	"github.com/orekhovskiy/shrtn/internal/logger"

	"github.com/stretchr/testify/assert"

	"github.com/orekhovskiy/shrtn/config"
)

func TestCreateShortUrl(t *testing.T) {
	mockLogger := &logger.NoopLogger{}
	mockService := new(mocks.MockURLService)
	opts := config.Config{BaseURL: "http://localhost:8080"}
	handler := Handler{logger: mockLogger, opts: opts, urlService: mockService}

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
			if tt.mockSaveReturn != "" {
				mockService.On("Save", tt.body).Return(tt.mockSaveReturn)
			}

			req := httptest.NewRequest(tt.method, "/shorten", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", tt.contentType)

			rec := httptest.NewRecorder()

			handler.CreateShortURL(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			responseBody, _ := io.ReadAll(rec.Body)
			assert.Equal(t, tt.expectedBody, string(responseBody))
		})
	}
}
