package mocks

import (
	"github.com/orekhovskiy/shrtn/internal/entity"
	"github.com/stretchr/testify/mock"
)

// Mocking the URL service
type MockURLService struct {
	mock.Mock
}

func (m *MockURLService) Save(originalURL string) (string, error) {
	args := m.Called(originalURL)
	return args.String(0), nil
}

func (m *MockURLService) GetByID(id string) (string, error) {
	args := m.Called(id)
	return args.String(0), args.Error(1)
}

func (m *MockURLService) Ping() error {
	return nil
}
func (m *MockURLService) ProcessBatch(batchRequests []entity.BatchRequest) ([]entity.BatchResponse, error) {
	args := m.Called(batchRequests)

	if resp, ok := args.Get(0).([]entity.BatchResponse); ok {
		return resp, args.Error(1)
	}
	return []entity.BatchResponse{}, args.Error(1)
}
