package repository

import (
	"github.com/stretchr/testify/mock"

	"go-films-api/internal/domain"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByUsername(username string) (*domain.User, error) {
	args := m.Called(username)
	if user, ok := args.Get(0).(*domain.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}
