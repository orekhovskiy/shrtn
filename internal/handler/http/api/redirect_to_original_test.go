package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"strings"
)

// Mocking the URL service
func TestRedirectToOriginal(t *testing.T) {
	mockService := new(MockURLService)
	handler := Handler{urlService: mockService}

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
