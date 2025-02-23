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

func TestRegister_Success(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	service := usecase.NewUserService(mockRepo)

	mockRepo.On("GetUserByUsername", "newuser").Return(nil, nil)
	mockRepo.On("CreateUser", mock.Anything).Return(nil)

	err := service.Register("newuser", "password123")
	assert.NoError(t, err)

	mockRepo.AssertCalled(t, "GetUserByUsername", "newuser")
	mockRepo.AssertCalled(t, "CreateUser", mock.Anything)
}

func TestRegister_UsernameTaken(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	service := usecase.NewUserService(mockRepo)

	existingUser := &domain.User{ID: 1, Username: "AlphaUser"}
	mockRepo.On("GetUserByUsername", "AlphaUser").Return(existingUser, nil)

	err := service.Register("AlphaUser", "secret12")
	assert.Error(t, err)
	assert.Equal(t, "username already taken", err.Error())
}

func TestRegister_InvalidUsername(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	service := usecase.NewUserService(mockRepo)

	err := service.Register("123Invalid", "somepass")
	assert.Error(t, err)
	assert.Equal(t, "username must start with a letter and contain only alphanumeric characters", err.Error())

	err = service.Register("John_Doe", "somepass")
	assert.Error(t, err)
	assert.Equal(t, "username must start with a letter and contain only alphanumeric characters", err.Error())
}

func TestRegister_PasswordTooShort(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	service := usecase.NewUserService(mockRepo)

	err := service.Register("AlphaUser", "123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "password must be between 6 and 20 characters")
}

func TestRegister_PasswordTooLong(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	service := usecase.NewUserService(mockRepo)

	tooLongPass := "thispasswordisdefinitelymorethan20chars"
	err := service.Register("BetaUser", tooLongPass)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "password must be between 6 and 20 characters")
}

func TestRegister_MissingUppercase(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	service := usecase.NewUserService(mockRepo)

	err := service.Register("UserTest", "abcd123#")
	assert.Error(t, err)
	assert.Equal(t, "password must contain at least one uppercase letter", err.Error())
}

func TestRegister_MissingDigit(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	service := usecase.NewUserService(mockRepo)

	err := service.Register("UserTest", "Abcd#xyz")
	assert.Error(t, err)
	assert.Equal(t, "password must contain at least one digit", err.Error())
}

func TestRegister_MissingSpecialChar(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	service := usecase.NewUserService(mockRepo)

	err := service.Register("UserTest", "Abcd1234")
	assert.Error(t, err)
	assert.Equal(t, "password must contain at least one special character", err.Error())
}

func TestRegister_ValidAllRequirements(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	service := usecase.NewUserService(mockRepo)

	validPassword := "Abcd1234!"
	mockRepo.On("GetUserByUsername", "ValidUser").Return(nil, nil)
	mockRepo.On("CreateUser", mock.Anything).Return(nil)

	err := service.Register("ValidUser", validPassword)
	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "CreateUser", mock.Anything)
}

func TestLogin_Success(t *testing.T) {
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

func TestLogin_InvalidPassword(t *testing.T) {
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

func TestLogin_NoUser(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	service := usecase.NewUserService(mockRepo)

	mockRepo.On("GetUserByUsername", "unknown").Return(nil, nil)

	token, exp, err := service.Login("unknown", "secret")
	assert.Empty(t, token)
	assert.Equal(t, time.Time{}, exp)
	assert.EqualError(t, err, "invalid username or password")
}
