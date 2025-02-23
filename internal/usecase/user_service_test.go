package usecase_test

import (
	"testing"
	"time"

	"go-films-api/internal/domain"
	"go-films-api/internal/repository"
	"go-films-api/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterSuccess(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	service := usecase.NewUserService(mockRepo)

	mockRepo.On("GetUserByUsername", "newuser").Return(nil, nil)
	mockRepo.On("CreateUser", mock.Anything).Return(nil)

	err := service.Register("newuser", "password123")
	assert.NoError(t, err)

	mockRepo.AssertCalled(t, "GetUserByUsername", "newuser")
	mockRepo.AssertCalled(t, "CreateUser", mock.Anything)
}

func TestRegisterDuplicateUser(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	service := usecase.NewUserService(mockRepo)

	existing := &domain.User{ID: 1, Username: "existinguser", Password: "hashedPass"}
	mockRepo.On("GetUserByUsername", "existinguser").Return(existing, nil)

	err := service.Register("existinguser", "somepass")
	assert.Error(t, err)
	assert.Equal(t, "username already taken", err.Error())
}

func TestLoginSuccess(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	service := usecase.NewUserService(mockRepo)

	// Provide a hashed password that will pass bcrypt check:
	hashed := "$2a$10$1fybhpdIC527ODopk5/FLu5L5o60g.2p1NGd7Zso75iv.R4siZm3e"
	user := &domain.User{ID: 42, Username: "johndoe", Password: hashed}

	mockRepo.On("GetUserByUsername", "johndoe").Return(user, nil)

	token, exp, err := service.Login("johndoe", "secret")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.WithinDuration(t, time.Now().Add(time.Hour), exp, 2*time.Second)
}

func TestLoginInvalidPassword(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	service := usecase.NewUserService(mockRepo)

	// user with a known hashed password
	hashed := "$2a$10$1fybhpdIC527ODopk5/FLu5L5o60g.2p1NGd7Zso75iv.R4siZm3e"
	user := &domain.User{ID: 42, Username: "johndoe", Password: hashed}

	mockRepo.On("GetUserByUsername", "johndoe").Return(user, nil)

	token, exp, err := service.Login("johndoe", "wrongpass")
	assert.Empty(t, token)
	assert.Equal(t, time.Time{}, exp)
	assert.EqualError(t, err, "invalid username or password")
}

func TestLoginNoUser(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	service := usecase.NewUserService(mockRepo)

	mockRepo.On("GetUserByUsername", "unknown").Return(nil, nil)

	token, exp, err := service.Login("unknown", "secret")
	assert.Empty(t, token)
	assert.Equal(t, time.Time{}, exp)
	assert.EqualError(t, err, "invalid username or password")
}
