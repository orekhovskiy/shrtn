package mocks

import "github.com/stretchr/testify/mock"

// Mocking the URL service
type MockURLService struct {
	mock.Mock
}

func (m *MockURLService) Save(originalURL string) string {
	args := m.Called(originalURL)
	return args.String(0)
}

func (m *MockURLService) GetByID(id string) (string, error) {
	args := m.Called(id)
	return args.String(0), args.Error(1)
}
