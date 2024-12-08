package userurls

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/orekhovskiy/shrtn/internal/handler/http/api/mocks"
	"github.com/orekhovskiy/shrtn/internal/service/authservice"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestDeleteUrls(t *testing.T) {
	logger := zap.NewNop() // Use a no-op logger for testing

	t.Run("unauthorized user", func(t *testing.T) {
		mockURLService := new(mocks.MockURLService)
		handler := Handler{logger: logger, urlService: mockURLService}

		req := httptest.NewRequest(http.MethodDelete, "/urls", nil)
		w := httptest.NewRecorder()

		handler.deleteUrls(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(t, "Unauthorized\n", w.Body.String())
	})

	t.Run("invalid request body", func(t *testing.T) {
		mockURLService := new(mocks.MockURLService)
		handler := Handler{logger: logger, urlService: mockURLService}

		ctx := context.WithValue(context.Background(), authservice.UserIDKey, "user-123")
		req := httptest.NewRequest(http.MethodDelete, "/urls", bytes.NewBuffer([]byte("invalid json")))
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.deleteUrls(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "Invalid request body\n", w.Body.String())
	})

	t.Run("successful deletion", func(t *testing.T) {
		mockURLService := new(mocks.MockURLService)
		handler := Handler{logger: logger, urlService: mockURLService}

		ctx := context.WithValue(context.Background(), authservice.UserIDKey, "user-123")
		shortURLs := []string{"short1", "short2"}
		body, _ := json.Marshal(shortURLs)
		req := httptest.NewRequest(http.MethodDelete, "/urls", bytes.NewBuffer(body))
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockURLService.On("MarkURLsAsDeleted", shortURLs, "user-123").Return(nil)

		handler.deleteUrls(w, req)

		assert.Equal(t, http.StatusAccepted, w.Code)
		mockURLService.AssertCalled(t, "MarkURLsAsDeleted", shortURLs, "user-123")
	})

	t.Run("error during deletion", func(t *testing.T) {
		mockURLService := new(mocks.MockURLService)
		handler := Handler{logger: logger, urlService: mockURLService}

		ctx := context.WithValue(context.Background(), authservice.UserIDKey, "user-123")
		shortURLs := []string{"short1", "short2"}
		body, _ := json.Marshal(shortURLs)
		req := httptest.NewRequest(http.MethodDelete, "/urls", bytes.NewBuffer(body))
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockURLService.On("MarkURLsAsDeleted", shortURLs, "user-123").Return([]error{errors.New("error 1"), errors.New("error 2")})

		handler.deleteUrls(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code, "expected status 500 Internal Server Error")
		assert.Equal(t, "Internal Server Error\n", w.Body.String(), "unexpected response body")
		mockURLService.AssertCalled(t, "MarkURLsAsDeleted", shortURLs, "user-123")
	})
}
