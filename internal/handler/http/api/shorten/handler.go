package shorten

import (
	"context"
	"github.com/orekhovskiy/shrtn/config"
	"github.com/orekhovskiy/shrtn/internal/entity"
	"github.com/orekhovskiy/shrtn/internal/logger"
)

type URLShortenerService interface {
	Save(url string, userID string) (id string, err error)
	ProcessBatch(batchRequests []entity.BatchRequest, userID string) (batchResponses []entity.BatchResponse, err error)
	BuildURL(uri string) string
}

type AuthService interface {
	GetUserIDFromContext(ctx context.Context) (string, bool)
}

type Handler struct {
	opts        config.Config
	urlService  URLShortenerService
	authService AuthService
	logger      logger.Logger
}

func NewHandler(logger logger.Logger, opts *config.Config, urlService URLShortenerService, authService AuthService) *Handler {
	return &Handler{logger: logger, opts: *opts, urlService: urlService, authService: authService}
}
