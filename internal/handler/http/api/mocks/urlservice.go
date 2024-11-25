package mocks

import (
	"github.com/orekhovskiy/shrtn/internal/entity"
	"github.com/stretchr/testify/mock"
)

// Mocking the URL service
type MockURLService struct {
	mock.Mock
}

func (m *MockURLService) Save(originalURL string, userID string) (string, error) {
	args := m.Called(originalURL, userID)
	return args.String(0), nil
}

func (m *MockURLService) GetByID(id string) (string, error) {
	args := m.Called(id)
	return args.String(0), args.Error(1)
}

func (m *MockURLService) Ping() error {
	return nil
}

func (m *MockURLService) ProcessBatch(batchRequests []entity.BatchRequest, userID string) ([]entity.BatchResponse, error) {
	args := m.Called(batchRequests, userID)

	if resp, ok := args.Get(0).([]entity.BatchResponse); ok {
		return resp, args.Error(1)
	}
	return []entity.BatchResponse{}, args.Error(1)
}

func (m *MockURLService) BuildURL(uri string) string {
	args := m.Called(uri)
	return args.String(0)
}

func (m *MockURLService) GetUserURLs(userID string) ([]entity.URLRecord, error) {
	args := m.Called(userID)
	if urls, ok := args.Get(0).([]entity.URLRecord); ok {
		return urls, args.Error(1)
	}
	return nil, args.Error(1)
}
