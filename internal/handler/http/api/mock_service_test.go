package api

import "github.com/stretchr/testify/mock"

// Mocking the URL service
type MockUrlService struct {
	mock.Mock
}

func (m *MockUrlService) Save(originalURL string) string {
	args := m.Called(originalURL)
	return args.String(0)
}

func (m *MockUrlService) GetById(id string) (string, error) {
	args := m.Called(id)
	return args.String(0), args.Error(1)
}
