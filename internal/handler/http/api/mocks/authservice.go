package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) GetUserIDFromContext(ctx context.Context) (string, bool) {
	args := m.Called(ctx)
	return args.String(0), args.Bool(1)
}
