package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/orekhovskiy/shrtn/internal/handler/http/api/mocks"

	"github.com/orekhovskiy/shrtn/config"

	"strings"

	"github.com/stretchr/testify/assert"
)

// Mocking the URL service
func TestRedirectToOriginal(t *testing.T) {
	mockLogger := &mocks.NoopLogger{}
	mockService := new(mocks.MockURLService)
	opts := config.Config{BaseURL: "http://localhost:8080"}
	handler := Handler{logger: mockLogger, opts: opts, urlService: mockService}

	tests := []struct {
		name             string
		path             string
		mockReturnURL    string
		mockReturnError  error
		expectedStatus   int
		expectedLocation string
	}{
		{
			name:             "Successful Redirect",
			path:             "/12345",
			mockReturnURL:    "http://example.com",
			mockReturnError:  nil,
			expectedStatus:   http.StatusTemporaryRedirect,
			expectedLocation: "http://example.com",
		},
		{
			name:             "URL Not Found",
			path:             "/invalid-id",
			mockReturnURL:    "",
			mockReturnError:  fmt.Errorf("URL not found"),
			expectedStatus:   http.StatusBadRequest,
			expectedLocation: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.On("GetByID", strings.TrimPrefix(tt.path, "/")).Return(tt.mockReturnURL, tt.mockReturnError)

			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()

			handler.RedirectToOriginal(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedStatus == http.StatusTemporaryRedirect {
				assert.Equal(t, tt.expectedLocation, rec.Header().Get("Location"))
			}
		})
	}
}
